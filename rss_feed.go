package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/prathamanvekar/gator/internal/database"
)

// RSSFeed represents the root structure of an RSS XML feed.
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents an individual post or item within an RSS feed.
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// fetchFeed fetches an RSS feed from a given URL, unmarshals the XML response
// into an RSSFeed structure, and decodes HTML entities in the titles and descriptions.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
}

// scrapeFeeds is a background worker function that gets the next feed that
// needs to be fetched, marks it as fetched in the database, fetches the content,
// and saves all items as posts to the database while ignoring duplicates.
func scrapeFeeds(s *state) {

	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("error getting the next feed to fetch: %v\n", err)
		return
	}

	markedFeed, err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		fmt.Printf("error marking the feed as fetched: %v\n", err)
		return
	}

	feed, err := fetchFeed(context.Background(), markedFeed.Url)
	if err != nil {
		fmt.Printf("error fetching feed: %v\n", err)
		return
	}

	for _, v := range feed.Channel.Item {

		newPostID := uuid.New()

		parsedTime, err := time.Parse(time.RFC1123Z, v.PubDate)
		var publishedAt sql.NullTime
		if err != nil {
			log.Printf("error parsing date: %v", err)
			// To insert a NULL value:
			publishedAt = sql.NullTime{
				Valid: false, // This triggers the NULL in SQL
			}
		} else {
			// To insert an actual Time:
			publishedAt = sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}
		}

		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          newPostID,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       v.Title,
			Url:         v.Link,
			Description: v.Description,
			PublishedAt: publishedAt,
			FeedID:      markedFeed.ID,
		})
		if err != nil {
			// Silently skip duplicate posts
			if err == sql.ErrNoRows || strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE") {
				continue
			}
			log.Printf("error creating post: %v", err)
		}

	}

}
