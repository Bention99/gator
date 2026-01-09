package main

import (
    "errors"
    "github.com/Bention99/gator/internal/config"
    "github.com/Bention99/gator/internal/database"
)

type state struct {
    db  *database.Queries
    cfg *config.Config
}

type command struct {
    name string
    args []string
}

type commands struct {
    cmds map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
    f, ok := c.cmds[cmd.name]
    if !ok {
        return errors.New("Command not supported")
    }
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmds[name] = f
}