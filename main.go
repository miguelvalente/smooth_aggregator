package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

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

	apiCfg := apiConfig{
		DB: *&dbQueries,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerHealthz)
	mux.HandleFunc("GET /v1/err", handlerErr)
	mux.HandleFunc("POST /v1/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /v1/users", apiCfg.middlewareAuth(apiCfg.handlerUsersGet))
	mux.HandleFunc("POST /v1/feeds", apiCfg.middlewareAuth(apiCfg.handlerFeedsCreate))
	mux.HandleFunc("GET /v1/feeds", apiCfg.handlerFeedsGet)
	mux.HandleFunc("POST /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedsFollowsCreate))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerFeedsFollowsDelete))
	mux.HandleFunc("GET /v1/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerFeedsFollowsGet))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScrapping(collectionInterval, apiCfg)

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

	//

	// feed_uuid := uuid.MustParse("aea1e783-bf58-48db-bbc4-0bb33ed3daab")
	// params := database.MarkFeedFetchedParams{
	// 	ID: feed_uuid,
	// 	LastFetchedAt: sql.NullTime{
	// 		Time:  time.Now().UTC(),
	// 		Valid: true,
	// 	},
	// }
	// _ = apiCfg.DB.MarkFeedFetched(context.TODO(), params)
	// what, _ := apiCfg.DB.GetNextNFeedsToFetch(context.TODO(), 10)
	// fmt.Println(what)

}
