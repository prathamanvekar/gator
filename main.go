package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/prathamanvekar/gator/internal/config"
	"github.com/prathamanvekar/gator/internal/database"
)

// state holds the application configuration and database queries
// that are passed around to the command handlers.
type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	// Read the configuration file to get the database URL and current user.
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	dbURL := cfg.DBURL

	// Open a connection to the PostgreSQL database.
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	defer db.Close()

	// Verify database connection works by pinging it.
	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	// Initialize the generated sqlc queries with the db connection.
	dbQueries := database.New(db)

	// Set up the application state.
	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	// Initialize the commands registry.
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register all available CLI commands and their corresponding handlers.
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	// Ensure at least a command name is provided.
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	// Parse the command name and arguments from the command line.
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// Run the requested command.
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
