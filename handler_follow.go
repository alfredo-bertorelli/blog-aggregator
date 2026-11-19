package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <URL>", cmd.Name)
	}
	url := cmd.Args[0]

	feedRow, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	feedFollows := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feedRow.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	feedFollowsRow, err := s.db.CreateFeedFollow(context.Background(), feedFollows)
	if err != nil {
		return fmt.Errorf("could't insert into feed_follows: %w", err)
	}
	fmt.Printf("User %s now follows feed %s!\n", feedFollowsRow.UserName, feedFollowsRow.FeedName)

	return nil
}
