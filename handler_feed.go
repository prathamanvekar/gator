package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prathamanvekar/gator/internal/database"
)

// handlerAddFeed creates a new RSS feed in the database and automatically
// follows it for the currently logged-in user.
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	feedID := uuid.New()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	followRecordID := uuid.New()
	followRecord, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        followRecordID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow record: %v", err)
	}

	fmt.Print("Followed Too!\n")
	fmt.Printf("Feed: %s\n", followRecord.FeedName)
	fmt.Printf("User: %s\n", followRecord.UserName)

	fmt.Print("Added!\n")
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("Url: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)

	return nil
}

// handlerFeeds lists all RSS feeds available in the database,
// displaying their name, URL, and the user who originally added them.
func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds from db: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, v := range feeds {
		fmt.Printf("Name: %v\n", v.Name)
		fmt.Printf("Url: %v\n", v.Url)
		user, err := s.db.GetUserByID(context.Background(), v.UserID)
		if err != nil {
			return fmt.Errorf("error getting user by id: %v", err)
		}
		fmt.Printf("User Name: %v\n", user.Name)
	}

	return nil
}
