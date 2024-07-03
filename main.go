package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

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

	bytes, err := ioutil.ReadAll(resp.Body)
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

func main() {
	as, _ := fetchRSS("https://wagslane.dev/index.xml")
	fmt.Println(as)
}

// func main() {
// 	godotenv.Load(".env")

// 	const filepathRoot = "."
// 	port := os.Getenv("PORT")
// 	dbURL := os.Getenv("DB_URL")

// 	db, _ := sql.Open("postgres", dbURL)
// 	dbQueries := database.New(db)

// 	apiConfig := apiConfig{
// 		DB: *&dbQueries,
// 	}

// 	mux := http.NewServeMux()

// 	mux.HandleFunc("GET /v1/healthz", handlerHealthz)
// 	mux.HandleFunc("GET /v1/err", handlerErr)
// 	mux.HandleFunc("POST /v1/users", apiConfig.handlerUsersCreate)
// 	mux.HandleFunc("GET /v1/users", apiConfig.middlewareAuth(apiConfig.handlerUsersGet))
// 	mux.HandleFunc("POST /v1/feeds", apiConfig.middlewareAuth(apiConfig.handlerFeedsCreate))
// 	mux.HandleFunc("GET /v1/feeds", apiConfig.handlerFeedsGet)
// 	mux.HandleFunc("POST /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerFeedsFollowsCreate))
// 	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.handlerFeedsFollowsDelete))
// 	mux.HandleFunc("GET /v1/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerFeedsFollowsGet))

// 	//

// 	feed_uuid := uuid.MustParse("aea1e783-bf58-48db-bbc4-0bb33ed3daab")
// 	params := database.MarkFeedFetchedParams{
// 		ID: feed_uuid,
// 		LastFetchedAt: sql.NullTime{
// 			Time:  time.Now().UTC(),
// 			Valid: true,
// 		},
// 	}
// 	_ = apiConfig.DB.MarkFeedFetched(context.TODO(), params)
// 	what, _ := apiConfig.DB.GetNextNFeedsToFetch(context.TODO(), 10)
// 	fmt.Println(what)
// 	srv := &http.Server{
// 		Addr:    ":" + port,
// 		Handler: mux,
// 	}

// 	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
// 	log.Fatal(srv.ListenAndServe())

// }
