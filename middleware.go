package main

import (
	"context"
	"fmt"

	"github.com/prathamanvekar/gator/internal/database"
)

// middlewareLoggedIn is a middleware function that ensures a user is logged in
// before allowing the wrapped command handler to execute. It fetches the
// currently configured user from the database and passes it to the handler.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error getting user: %v", err)
		}
		return handler(s, cmd, user)
	}
}
