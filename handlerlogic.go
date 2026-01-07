package main

import (
	"errors"
	"fmt"
	"github.com/Bention99/gator/internal/config"
)

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("login expects a single argument, the username")
    }
    s.cfg.DBURL = "postgres://example"
	s.cfg.CurrentUserName = cmd.args[0]

	err := config.SetUser(*s.cfg)
    if err != nil {
		return err
	}

    fmt.Printf("User set to: %v\n", cmd.args[0])
    return nil
}