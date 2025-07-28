package tables

import (
	"log"

	"github.com/gocql/gocql"
)

func CreateTasksTable(session *gocql.Session) {
	// Create tasks table
	query := `
        CREATE TABLE IF NOT EXISTS tasks (
            task_id UUID,
            user_id UUID,
            title TEXT,
            description TEXT,
            status TEXT,
            created_at TIMESTAMP,
            updated_at TIMESTAMP,
            PRIMARY KEY (task_id)
        );
    `
	if err := session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create 'tasks' table: %v", err)
	}

	// Create index on user_id
	indexQuery := `CREATE INDEX IF NOT EXISTS ON tasks (user_id);`
	if err := session.Query(indexQuery).Exec(); err != nil {
		log.Fatalf("Failed to create index on tasks.user_id: %v", err)
	}

	log.Println("'tasks' table and indices created successfully!")
}
