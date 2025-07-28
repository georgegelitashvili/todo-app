package models

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

type UserNotFoundError struct {
	Email string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found with email: %s", e.Email)
}

type User struct {
	UserID    gocql.UUID `json:"user_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password,omitempty"` // Don't include in JSON responses
	CreatedAt time.Time  `json:"created_at"`
}

func NewUser(username, email, password string) *User {
	return &User{
		UserID:    gocql.TimeUUID(), // Generate proper UUID
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now().UTC(),
	}
}

func (u *User) Create(session *gocql.Session) error {
	// Set UUID if not set
	if u.UserID == (gocql.UUID{}) {
		u.UserID = gocql.TimeUUID()
	}

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("password hashing error: %v", err)
	}

	// Debug log
	log.Printf("Creating user with UUID: %s", u.UserID)

	query := `INSERT INTO users (user_id, username, email, password, created_at) 
             VALUES (?, ?, ?, ?, ?)`

	if err := session.Query(query,
		u.UserID,
		u.Username,
		u.Email,
		hashedPassword,
		u.CreatedAt).Exec(); err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	return nil
}

func GetUserByEmail(session *gocql.Session, email string) (*User, error) {
	user := &User{}
	query := `SELECT user_id, username, email, password, created_at 
             FROM users WHERE email = ? ALLOW FILTERING`

	err := session.Query(query, email).
		Consistency(gocql.One).
		Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == gocql.ErrNotFound {
		log.Printf("No user found with email: %s", email)
		return nil, &UserNotFoundError{Email: email}
	} else if err != nil {
		log.Printf("Query error: %v", err)
		return nil, fmt.Errorf("query error: %v", err)
	}

	// Verify user data is valid
	emptyUUID := gocql.UUID{}
	if user.UserID == emptyUUID {
		log.Printf("Invalid UUID for user: %s", email)
		return nil, fmt.Errorf("invalid user data retrieved: empty UUID")
	}

	return user, nil
}

func GetUserByID(session *gocql.Session, userID gocql.UUID) (*User, error) {
	user := &User{}
	query := `SELECT user_id, username, email, password, created_at FROM users WHERE user_id = ?`

	err := session.Query(query, userID).Consistency(gocql.One).
		Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err == gocql.ErrNotFound {
		return nil, fmt.Errorf("user not found with ID: %s", userID)
	} else if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}

	return user, nil
}

func DeleteUserByID(session *gocql.Session, userID gocql.UUID) error {
	query := `DELETE FROM users WHERE user_id = ?`
	return session.Query(query, userID).Exec()
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
