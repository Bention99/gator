package main

import (
	"fmt"
	"os"
	"github.com/Bention99/gator/internal/config"
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

	s := &state{
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