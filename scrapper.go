package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/miguelvalente/smooth_aggregator/internal/database"
)

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

func startScrapping(interval time.Duration, concurrency int, apiCfg apiConfig) {
	for {
		feeds, err := apiCfg.DB.GetNextNFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}

		for _, feed := range feeds {
			rss, _ := fetchRSS(feed.Url)
			err = apiCfg.DB.MarkFeedFetched(context.Background(), feed.ID)

			for _, post := range rss.Channel.Items {
				fmt.Println(post.Title)

				postParams := database.CreatePostParams{
					ID:        uuid.New(),
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
					Title: sql.NullString{
						String: post.Title,
						Valid:  true,
					},
					Description: sql.NullString{
						String: post.Description,
						Valid:  true,
					},
					Url: post.Link,
					PublishedAt: sql.NullTime{
						Time:  parseDate(post.PubDate),
						Valid: true,
					},
					FeedID: feed.ID,
				}
				apiCfg.DB.CreatePost(context.Background(), postParams)

			}

		}
	}
}

func parseDate(dateStr string) time.Time {
	layout := "2006-01-02T15:04:05Z" // Adjust the layout according to your date format
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{} // Return zero time on error
	}
	return t
}
