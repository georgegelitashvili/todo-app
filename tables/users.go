package tables

import (
	"log"

	"github.com/gocql/gocql"
)

// CreateUserTable creates a table for users if it doesn't exist

func CreateUsersTable(session *gocql.Session) {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			email TEXT,
			user_id UUID,
			username TEXT,
			password TEXT,
			created_at TIMESTAMP,
			PRIMARY KEY (email, user_id)
		);
	`
	if err := session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create 'users' table: %v", err)
	}
	log.Println("'users' table created successfully!")
}
