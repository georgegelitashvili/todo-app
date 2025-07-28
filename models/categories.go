package models

import (
	"time"

	"github.com/gocql/gocql"
)

type Category struct {
	CategoryID gocql.UUID `json:"category_id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Create method
func (c *Category) Create(session *gocql.Session) error {
	c.CategoryID = gocql.TimeUUID()
	c.CreatedAt = time.Now()
	query := `INSERT INTO categories (category_id, name, created_at) VALUES (?, ?, ?)`
	return session.Query(query, c.CategoryID, c.Name, c.CreatedAt).Exec()
}

// Update method
func (c *Category) Update(session *gocql.Session) error {
	query := `UPDATE categories SET name = ? WHERE category_id = ?`
	return session.Query(query, c.Name, c.CategoryID).Exec()
}

// Get methods
func GetCategoryByID(session *gocql.Session, categoryID gocql.UUID) (*Category, error) {
	category := &Category{}
	query := `SELECT category_id, name, created_at FROM categories WHERE category_id = ?`
	err := session.Query(query, categoryID).Scan(
		&category.CategoryID,
		&category.Name,
		&category.CreatedAt)
	return category, err
}

func GetAllCategories(session *gocql.Session) ([]Category, error) {
	var categories []Category
	query := "SELECT category_id, name, created_at FROM categories"
	iter := session.Query(query).Iter()
	var category Category
	for iter.Scan(
		&category.CategoryID,
		&category.Name,
		&category.CreatedAt) {
		categories = append(categories, category)
	}
	return categories, iter.Close()
}

// Delete method
func DeleteCategoryByID(session *gocql.Session, categoryID gocql.UUID) error {
	query := `DELETE FROM categories WHERE category_id = ?`
	return session.Query(query, categoryID).Exec()
}
