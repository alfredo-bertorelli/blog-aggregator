package main

import (
	"context"
	"fmt"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {

	innerFunc := func(s *state, cmd command) error {
		userRow, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("couldn't fetch user: %w", err)
		}
		return handler(s, cmd, userRow)
	}
	return innerFunc

}
