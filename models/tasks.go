package models

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

const (
	StatusPending    = "todo"
	StatusInProgress = "in_progress"
	StatusCompleted  = "done"
)

type Task struct {
	TaskID      gocql.UUID `json:"task_id"`
	UserID      gocql.UUID `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func NewTask(userID gocql.UUID, title, description, status string) *Task {
	if !isValidStatus(status) {
		status = StatusPending
	}
	return &Task{
		TaskID:      gocql.TimeUUID(), // Generate unique TimeUUID for each task
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      status,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (t *Task) Create(session *gocql.Session) error {
	query := `INSERT INTO tasks (task_id, user_id, title, description, status, created_at, updated_at) 
             VALUES (?, ?, ?, ?, ?, ?, ?)`

	return session.Query(query,
		t.TaskID,
		t.UserID,
		t.Title,
		t.Description,
		t.Status,
		t.CreatedAt,
		t.UpdatedAt).Exec()
}

func GetTasksByUserID(session *gocql.Session, userID gocql.UUID) ([]*Task, error) {
	var tasks []*Task

	query := `SELECT task_id, user_id, title, description, status, created_at, updated_at 
             FROM tasks WHERE user_id = ? ALLOW FILTERING`

	iter := session.Query(query, userID).Iter()

	var task Task
	for iter.Scan(
		&task.TaskID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt) {
		tasks = append(tasks, &Task{
			TaskID:      task.TaskID,
			UserID:      task.UserID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *Task) Update(session *gocql.Session) error {
	if !isValidStatus(t.Status) {
		return fmt.Errorf("invalid status: %s", t.Status)
	}
	t.UpdatedAt = time.Now()
	query := `UPDATE tasks 
			 SET title = ?, description = ?, status = ?, updated_at = ? 
			 WHERE task_id = ?`
	return session.Query(query,
		t.Title,
		t.Description,
		t.Status,
		t.UpdatedAt,
		t.TaskID).Exec()
}

func DeleteTaskByID(session *gocql.Session, taskID gocql.UUID) error {
	query := `DELETE FROM tasks WHERE task_id = ?`
	return session.Query(query, taskID).Exec()
}

func isValidStatus(status string) bool {
	return status == StatusPending ||
		status == StatusInProgress ||
		status == StatusCompleted
}
