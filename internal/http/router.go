package http

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	h := NewHandler(db)

	r.Use(middleware.Recoverer)
	r.Use(RequestID)
	r.Use(RequestLogger)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Endpoint not found"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", HealthHandler)
		r.Get("/ping", PingHandler)
		r.Get("/version", VersionHandler)

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", h.ListTasksHandler)
			r.Post("/", h.CreateTaskHandler)
			r.Get("/{id}", h.GetTaskHandler)
			r.Put("/{id}", h.UpdateTaskHandler)
			r.Delete("/{id}", h.DeleteTaskHandler)
		})
	})

	return r
}
