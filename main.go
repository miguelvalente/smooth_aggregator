package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
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

// Item represents an individual item in the RSS feed
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Channel represents the channel element in RSS
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Channel represents the channel element in RSS
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

func fetchRSS(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	err = xml.Unmarshal(bytes, &rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil

}

func smooth_worker(n int32, apiCfg apiConfig) {
	ctx := context.TODO()
	for {
		feedsToFech, _ := apiCfg.DB.GetNextNFeedsToFetch(ctx, n)
		for _, feed := range feedsToFech {
			rss, _ := fetchRSS(feed.Url)
			fmt.Println()
			fmt.Printf("=================")
			for i := 0; i < len(feed.Url); i++ {
				fmt.Print("=")
			}
			fmt.Printf("=================")
			fmt.Print("\n")

			fmt.Println("================", feed.Url, "================")
			fmt.Println()
			for _, post := range rss.Channel.Items {
				fmt.Println(post.Title)
			}

		}
		fmt.Println("\nSleepy time for sitty secons ZZZ")
		time.Sleep(60 * time.Second)
		fmt.Println("\nWhat its monin aweady. Wakupy timi")
		time.Sleep(2 * time.Second)
	}
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

	go smooth_worker(10, apiCfg)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

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
