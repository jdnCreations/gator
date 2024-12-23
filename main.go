package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jdnCreations/gator/internal/config"
	"github.com/jdnCreations/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
  db *database.Queries
	cfg *config.Config
}


func main() {
  
  
  cfg, err := config.Read()
	if err != nil {
    log.Fatal(err)
	}
  
  dbURL := cfg.DBUrl

  db, err := sql.Open("postgres", dbURL)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()
  db.Ping()

  dbQueries := database.New(db)

	programState := &state{
    db: dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
  cmds.register("register", handlerRegister)
  cmds.register("reset", handlerReset)
  cmds.register("users", handlerUsers)
  cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
  cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("usages: cli <command> [args...]")
		os.Exit(1)	
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
