package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/Bharat1Rajput/taskflow-backend/internal/repository"
)

func main() {
	port := os.Getenv("API_PORT")
	dbURL := os.Getenv("DATABASE_URL")

	if port == "" {
		port = "8080"
	}

	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	//  Run migrations BEFORE server starts
	if err := repository.RunMigrations(dbURL); err != nil {
		log.Fatal("Migration failed:", err)
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	log.Println("Server running on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
