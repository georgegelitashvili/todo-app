package keyspace

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

// createTodoKeyspace creates a keyspace for the todo app if it doesn't exist
func CreateTodoKeyspace(session *gocql.Session) {
	query := `
		CREATE KEYSPACE IF NOT EXISTS todo
		WITH REPLICATION = {
		'class' : 'SimpleStrategy',
		'replication_factor' : 3
	};`

	if err := session.Query(query).Exec(); err != nil {
		fmt.Println(err)
	}
	log.Println("Created keyspace todo")
}
