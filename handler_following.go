package main

import (
	"context"
	"fmt"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
)

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeedsByUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("couldn't fetch feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
	}

	return nil
}
