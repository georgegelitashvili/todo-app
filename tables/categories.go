package tables

import (
	"log"

	"github.com/gocql/gocql"
)

// CreateCategoriesTable creates the 'categories' table.
func CreateCategoriesTable(session *gocql.Session) {
	query := `
		CREATE TABLE IF NOT EXISTS categories (
			category_id UUID PRIMARY KEY,
			name TEXT,
			created_at TIMESTAMP
		);
	`
	if err := session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create 'categories' table: %v", err)
	}
	log.Println("'categories' table created successfully!")
}
