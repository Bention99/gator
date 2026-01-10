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

	cs := availableCommands()
	
	c := readConfig()

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

	name := args[1]
	cmdArgs := args[2:]
	
	cmd := command{
		name: name,
		args: cmdArgs,
	}

	if err := cs.run(s, cmd); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	/*c = readConfig()
	fmt.Printf("DBURL: %v User: %v\n", c.DBURL, c.CurrentUserName)*/
}

func availableCommands() commands {
	cs := commands{
		cmds: make(map[string]func(*state, command) error),
	}
	cs.register("login", handlerLogin)
	cs.register("register", handlerRegister)
	cs.register("reset", handlerReset)
	cs.register("users", handlerGetAllUsers)
	return cs
}

func readConfig() config.Config {
	c, err := config.Read()
	if err != nil {
		fmt.Println("Error reading File:")
		fmt.Printf(" - %v\n", err)
		os.Exit(1)
	}
	return c
}