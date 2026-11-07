package main

import (
	"log"
	"os"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	appState := &state{
		config: &cfg,
	}
	cmds := commands{
		registry: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	osArgs := os.Args
	if len(osArgs) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmd := command{
		Name: osArgs[1],
		Args: osArgs[2:],
	}
	err = cmds.run(appState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
