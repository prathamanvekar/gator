package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/prathamanvekar/gator/internal/database"
)

// handlerBrowse allows a logged-in user to view the latest posts from
// the feeds they follow. It takes an optional limit parameter to specify
// how many posts to retrieve (defaults to 2).
func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("error converting limits: %v", err)
		}
		if parsedLimit <= 0 {
			return fmt.Errorf("limit must be a positive integer")
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts for user: %v", err)
	}

	for _, v := range posts {
		fmt.Println("========================")

		fmt.Printf("Title: %s\n", v.Title)
		fmt.Printf("Url: %s\n", v.Url)
		fmt.Printf("Desc: %s\n", v.Description)

		fmt.Println("========================")
	}

	return nil
}
