package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/miguelvalente/smooth_aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	const filepathRoot = "."
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	apiConfig := apiConfig{
		DB: *&dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerHealthz)
	mux.HandleFunc("GET /v1/err", handlerErr)
	mux.HandleFunc("POST /v1/users", apiConfig.handlerUsersCreate)
	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.handlerUsersGet))
	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.handlerFeedsCreate))
	mux.HandleFunc("GET /v1/feeds", apiConfig.handlerFeedsGet)
	mux.HandleFunc("POST /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerFeedsFollowsCreate))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerFeedsFollowsDelete))
	mux.HandleFunc("GET /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerFeedsFollowsGet))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
