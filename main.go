package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

var (
	tasks  = []Task{}
	nextID = 1
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed\n"))
		return
	}

	w.Write([]byte("OK\n"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed\n"))
		return
	}

	w.Write([]byte("pong\n"))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed\n"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"version": "1.0.0"}` + "\n"))
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var newTask Task
		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
			return
		}
		newTask.ID = nextID
		nextID++
		tasks = append(tasks, newTask)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed\n"))
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/version", versionHandler)
	mux.HandleFunc("/tasks", tasksHandler)

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
