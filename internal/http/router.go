package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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
			r.Get("/", ListTasksHandler)
			r.Post("/", CreateTaskHandler)
			r.Get("/{id}", GetTaskHandler)
			r.Put("/{id}", UpdateTaskHandler)
			r.Delete("/{id}", DeleteTaskHandler)
		})
	})

	return r
}
