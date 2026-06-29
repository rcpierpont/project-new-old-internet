package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rcpierpont/project-new-old-internet/internal/database"
)

type envConfig struct {
	platform string
	db       *database.Queries
}

// TODO: db schema and basic queries for weeoos(posts) and honks(comments) and basic CRUD operations
func main() {
	godotenv.Load()

	const filepathRoot = "."
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	fmt.Printf("db url is: %s\n", dbURL)
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	cfg := envConfig{
		platform: platform,
		db:       database.New(db),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareHTTP(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("POST /api/users", cfg.handlerUsers)
	mux.HandleFunc("GET /api/users/{userID}", cfg.handlerGetUserByID)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *envConfig) middlewareHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
