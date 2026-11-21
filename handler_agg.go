package main

import (
	"fmt"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <interval>", cmd.Name)
	}
	interval, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't parse interval: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n\n", interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		scrapeFeed(s)
	}

	return nil
}
