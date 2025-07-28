package routes

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"todo-app/controllers"
	"todo-app/middleware"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type RouterConfig struct {
	Session       *gocql.Session
	Templates     *template.Template
	ComponentsDir string
}

func NewRouter(config RouterConfig) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.NoCacheMiddleware)

	// Static file server
	fs := http.FileServer(http.Dir("static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Serve component templates with version parameter
	router.HandleFunc("/templates/components/{component}.html", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		component := vars["component"]

		// Add version parameter to prevent caching
		w.Header().Set("ETag", time.Now().String())
		log.Printf("Loading component with no-cache: %s", component)

		http.ServeFile(w, r, filepath.Join(config.ComponentsDir, component+".html"))
	}).Methods("GET")

	// Controllers initialization
	userCtrl := controllers.NewUserController(config.Session)
	taskCtrl := controllers.NewTaskController(config.Session)
	categoryCtrl := controllers.NewCategoryController(config.Session)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/login", userCtrl.Login).Methods("POST")
	api.HandleFunc("/register", userCtrl.CreateUser).Methods("POST")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Protected User routes
	protected.HandleFunc("/users/{id}", userCtrl.GetUser).Methods("GET")
	protected.HandleFunc("/users/{id}", userCtrl.DeleteUser).Methods("DELETE")

	// Protected Task routes
	protected.HandleFunc("/tasks", taskCtrl.CreateTask).Methods("POST")
	protected.HandleFunc("/tasks/{id}", taskCtrl.GetTask).Methods("GET")
	protected.HandleFunc("/tasks", taskCtrl.GetAllTasks).Methods("GET")
	protected.HandleFunc("/tasks/{id}", taskCtrl.UpdateTask).Methods("PUT")
	protected.HandleFunc("/tasks/{id}", taskCtrl.DeleteTask).Methods("DELETE")

	// Protected Category routes
	protected.HandleFunc("/categories", categoryCtrl.CreateCategory).Methods("POST")
	protected.HandleFunc("/categories/{id}", categoryCtrl.GetCategory).Methods("GET")
	protected.HandleFunc("/categories", categoryCtrl.GetAllCategories).Methods("GET")
	protected.HandleFunc("/categories/{id}", categoryCtrl.UpdateCategory).Methods("PUT")
	protected.HandleFunc("/categories/{id}", categoryCtrl.DeleteCategory).Methods("DELETE")

	// Main route handler
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		config.Templates.ExecuteTemplate(w, "index.html", nil)
	}).Methods("GET")

	return router
}
