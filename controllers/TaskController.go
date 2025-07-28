package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-app/models"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

type TaskController struct {
	session *gocql.Session
}

func NewTaskController(session *gocql.Session) *TaskController {
	return &TaskController{session: session}
}

func (c *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := task.Create(c.session); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Status: "success",
		Data:   task,
	})
}

func (c *TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context
	userID, ok := r.Context().Value("user_id").(gocql.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get task_id from URL params
	params := mux.Vars(r)
	taskID, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get task and verify ownership
	task, err := models.GetTasksByUserID(c.session, taskID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	// Verify tasks belong to user
	for _, t := range task {
		if t.UserID != userID {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data:   task,
	})
}

func (c *TaskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(gocql.UUID)
	if !ok {
		log.Printf("Failed to get user_id from context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get tasks for specific user
	tasks, err := models.GetTasksByUserID(c.session, userID)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data: map[string]interface{}{
			"tasks": tasks,
		},
	})
}

func (c *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	task.TaskID = id
	if err := task.Update(c.session); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Task updated successfully",
		Data:    task,
	})
}

func (c *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Get user_id from context
	userID, ok := r.Context().Value("user_id").(gocql.UUID)
	if !ok {
		log.Printf("Failed to get user_id from context")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get task_id from URL
	params := mux.Vars(r)
	taskID, err := gocql.ParseUUID(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get task to verify ownership
	task, err := models.GetTasksByUserID(c.session, taskID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	// Verify tasks belong to user
	for _, t := range task {
		if t.UserID != userID {
			respondWithError(w, http.StatusForbidden, "Access denied")
			return
		}
	}

	// Delete task
	if err := models.DeleteTaskByID(c.session, taskID); err != nil {
		log.Printf("Failed to delete task: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Task deleted successfully",
	})
}
