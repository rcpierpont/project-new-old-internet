package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rcpierpont/project-new-old-internet/internal/database"
)

type envConfig struct {
	platform string
	db       *database.Queries
	hits     atomic.Int32
	secret   string
}

// TODO: authenticated queries for creating and fetching kreeyaws
// TODO: basic front end client that takes email and password and posts to login endpoint
func main() {
	godotenv.Load()

	const filepathRoot = "app"
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	cfg := envConfig{
		platform: platform,
		db:       database.New(db),
		hits:     atomic.Int32{},
		secret:   secret,
	}

	mux := http.NewServeMux()
	appHandler := cfg.middlewareHTTP(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", appHandler)

	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /api/users/{userID}", cfg.handlerGetUserByID)

	mux.HandleFunc("POST /api/kreeyaws", cfg.handlerCreateKreeyaw)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)

	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *envConfig) middlewareHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}
