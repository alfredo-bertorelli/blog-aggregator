package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	userCreated, err := s.db.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error registering username in DB: %s", err)
	}
	err = s.config.SetUser(userCreated.Name)
	if err != nil {
		return fmt.Errorf("error updating username in config file: %s", err)
	}
	fmt.Printf("New username %s registered!/\n", s.config.CurrentUserName)
	return nil
}
