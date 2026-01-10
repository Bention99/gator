package main

import (
	"errors"
	"fmt"
	"time"
	"database/sql"
	"context"
	"github.com/google/uuid"
	"github.com/Bention99/gator/internal/config"
	"github.com/Bention99/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("login expects a single argument, the username")
    }

	ctx := context.Background()

	user, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil {
		return errors.New("No Account for that username.")
	}

	s.cfg.CurrentUserName = cmd.args[0]

	err = config.SetUser(*s.cfg)
    if err != nil {
		return err
	}

    fmt.Printf("Login succesfull. Welcome %v\n", user.Name)
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("Register expects a single argument: Name")
    }

	ctx := context.Background()

	_, err := s.db.GetUser(ctx, cmd.args[0])
	if err == nil {
		return errors.New("User already exists.")
	}

	now := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	cup := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name: cmd.args[0],
	}

	user, err := s.db.CreateUser(ctx, cup)
	if err != nil {
		return err
	}
	s.cfg.CurrentUserName = user.Name

	err = config.SetUser(*s.cfg)
    if err != nil {
		return err
	}

    fmt.Printf("User created\n")
	fmt.Printf("ID: %v\n", user.ID)
	fmt.Printf("CreatedAt: %v\n", user.CreatedAt)
	fmt.Printf("updated_at: %v\n", user.UpdatedAt)
	fmt.Printf("Name: %v\n", user.Name)
    return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Println("Reset successful.")
	return nil
}

func handlerGetAllUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %v (current)\n", user)
		} else {
			fmt.Printf("* %v\n", user)
		}
	}
	return nil
}