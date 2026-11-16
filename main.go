package main

import (
	"log"
	"os"

	"database/sql"

	"github.com/alfredo-bertorelli/blog-aggregator/internal/config"
	"github.com/alfredo-bertorelli/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	appState := &state{
		config: &cfg,
		db:     database.New(db),
	}
	cmds := commands{
		registry: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(appState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
