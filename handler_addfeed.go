package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		return fmt.Errorf("couldn't add feed for user %s: %s", user.Name, err)
	}
	fmt.Printf("Feed row: %+v", feed)

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("couldn't add feed_follow record: %w", err)
	}

	return nil
}
