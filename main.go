package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"todo-app/config"
	"todo-app/keyspace"
	"todo-app/routes"
	"todo-app/tables"
)

func main() {
	// Database setup
	time.Sleep(5 * time.Second)
	systemSession := config.ConnectToCassandra("system")
	if systemSession == nil {
		log.Fatal("Failed to connect to Cassandra system keyspace")
	}

	keyspace.CreateTodoKeyspace(systemSession)
	systemSession.Close()

	time.Sleep(2 * time.Second)

	todoSession := config.ConnectToCassandra("todo")
	if todoSession == nil {
		log.Fatal("Failed to connect to todo keyspace")
	}
	defer todoSession.Close()

	// Create tables
	tables.CreateUsersTable(todoSession)
	tables.CreateTasksTable(todoSession)
	tables.CreateCategoriesTable(todoSession)

	// Initialize router from routes package
	workDir, _ := os.Getwd()
	templatesDir := filepath.Join(workDir, "templates")
	componentsDir := filepath.Join(templatesDir, "components")
	templates := template.Must(template.ParseGlob(filepath.Join(templatesDir, "*.html")))

	// Initialize router with config
	routerConfig := routes.RouterConfig{
		Session:       todoSession,
		Templates:     templates,
		ComponentsDir: componentsDir,
	}

	router := routes.NewRouter(routerConfig)

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
