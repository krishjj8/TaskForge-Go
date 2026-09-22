package main

import (
	"log"
	"net/http"

	"taskforge/internal/database"
	appHTTP "taskforge/internal/http"

	"taskforge/internal/config"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)

	}
	defer db.Close()
	log.Println("Database connection pool established successfully")
	router := appHTTP.NewRouter(db)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s in %s mode...\n", addr, cfg.Env)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
