package main

import (
	_ "github.com/lib/pq"
	"fmt"
	"os"
	"database/sql"
	"github.com/Bention99/gator/internal/config"
	"github.com/Bention99/gator/internal/database"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Please provide the command and 1 argument")
		os.Exit(1)
	}

	c, err := config.Read()
	if err != nil {
		fmt.Println("Error reading File:")
		fmt.Printf(" - %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", c.DBURL)
	if err != nil {
		fmt.Printf("Error connecting to the DB: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	s := &state{
		db: dbQueries,
		cfg: &c,
	}

	cs := commands{
		cmds: make(map[string]func(*state, command) error),
	}

	name := args[1]
	cmdArgs := args[2:]
	
	cmd := command{
		name: name,
		args: cmdArgs,
	}

	cs.register("login", handlerLogin)

	if err := cs.run(s, cmd); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	c2, err := config.Read()
	fmt.Printf("DBURL: %v User: %v\n", c2.DBURL, c2.CurrentUserName)
}