package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"taskforge/internal/model"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK\n"))
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong\n"))
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"version": "1.0.0"}` + "\n"))
}

// GET /api/v1/tasks
func (h *Handler) ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, title, COALESCE(description, ''), done FROM tasks ORDER BY id ASC`

	rows, err := h.db.QueryContext(r.Context(), query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch tasks"})
		return
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to scan task"})
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// POST /api/v1/tasks
func (h *Handler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask model.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	query := `
		INSERT INTO tasks (title, description, done)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := h.db.QueryRowContext(r.Context(), query, newTask.Title, newTask.Description, newTask.Done).Scan(&newTask.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// GET /api/v1/tasks/{id}
func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	query := `SELECT id, title, COALESCE(description, ''), done FROM tasks WHERE id = $1`

	var t model.Task
	err = h.db.QueryRowContext(r.Context(), query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Done)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// PUT /api/v1/tasks/{id}
func (h *Handler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	var updatedData model.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
		return
	}

	query := `
		UPDATE tasks
		SET title = $1, description = $2, done = $3
		WHERE id = $4
		RETURNING id
	`

	err = h.db.QueryRowContext(r.Context(), query, updatedData.Title, updatedData.Description, updatedData.Done, id).Scan(&updatedData.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update task"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedData)
}

// DELETE /api/v1/tasks/{id}
func (h *Handler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	query := `DELETE FROM tasks WHERE id = $1`

	result, err := h.db.ExecContext(r.Context(), query, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete task"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
