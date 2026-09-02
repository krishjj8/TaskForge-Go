package main

import (
	"log"
	"net/http"

	appHTTP "taskforge/internal/http"

	"taskforge/internal/config"
)

func main() {
	cfg := config.Load()

	router := appHTTP.NewRouter()

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s in %s mode...\n", addr, cfg.Env)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
