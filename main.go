package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/miguelvalente/smooth_aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB      *database.Queries
	Context context.Context
}

func main() {
	godotenv.Load(".env")

	const filepathRoot = "."
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	ctx := context.Background()

	apiConfig := apiConfig{
		DB:      *&dbQueries,
		Context: ctx,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerHealthz)
	mux.HandleFunc("GET /v1/err", handlerErr)
	mux.HandleFunc("POST /v1/users", apiConfig.handlerUsers)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
