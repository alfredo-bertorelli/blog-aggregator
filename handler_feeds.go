package main

import (
	"context"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <feed_name> <feed_url>", cmd.Name)
	}
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get user: %s", err)
	}

	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Username: %s\n", feed.Name_2.String)
		fmt.Println()
	}

	return nil
}
