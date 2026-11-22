package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s [limit]", cmd.Name)
	}
	if len(cmd.Args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil || limit < 1 {
			return fmt.Errorf("[limit] argument must be a positive integer")
		}
		limit = parsedLimit
	}

	queryArgs := database.GetRecentPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetRecentPostsForUser(context.Background(), queryArgs)
	if err != nil {
		return fmt.Errorf("couldn't fetch posts: %w", err)
	}

	for i, post := range posts {
		fmt.Printf("Post #%d: %s - %s\n", i+1, post.Title, post.Url)
	}

	return nil
}
