package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prathamanvekar/gator/internal/database"
)

// handlerFollow allows a user to follow an existing feed by providing its URL.
// It retrieves the feed by URL and creates a new follow record.
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	followRecordID := uuid.New()

	feed, err := s.db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}

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

	fmt.Printf("Followed %s!\n", cmd.Args[0])
	fmt.Printf("Feed: %s\n", followRecord.FeedName)
	fmt.Printf("User: %s\n", followRecord.UserName)

	return nil
}

// handlerFollowing displays a list of all feeds the currently logged-in
// user is following.
func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	follow_feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting follow feeds: %v", err)
	}
	fmt.Print("Feeds:\n")
	for _, v := range follow_feeds {
		fmt.Printf("%s\n", v.FeedName)
	}

	return nil
}

// handlerUnfollow removes a feed follow record for the currently logged-in user,
// allowing them to stop following a specific feed by its URL.
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	feed, err := s.db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}

	err = s.db.DeleteFeedFollowRecord(context.Background(), database.DeleteFeedFollowRecordParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error deleting feed follow record: %v", err)
	}

	fmt.Print("Record Deleted!\n")
	fmt.Printf("Unfollowed %s!\n", cmd.Args[0])

	return nil
}
