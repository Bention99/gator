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
	"github.com/Bention99/gator/internal/rss"
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

    fmt.Printf("Login succesfull. Welcome, %v!\n", user.Name)
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
	fmt.Printf("UpdatedAt: %v\n", user.UpdatedAt)
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

func handlerAggregate(s *state, cmd command) error {
	ctx := context.Background()
	feed, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%v", feed.Channel.Title)
	fmt.Printf("%v", feed.Channel.Description)
	for i := range feed.Channel.Item {
		fmt.Printf("%v", feed.Channel.Item[i].Title)
		fmt.Printf("%v", feed.Channel.Item[i].Description)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	ctx := context.Background()
	currentUserName := s.cfg.CurrentUserName
	currentUser, err := s.db.GetUser(ctx, currentUserName)
	if err != nil {
		return err
	}

	now := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	cfp := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name: cmd.args[0],
		Url: cmd.args[1],
		UserID: currentUser.ID,
	}

	feed, err := s.db.CreateFeed(ctx, cfp)
	if err != nil {
		return err
	}
	fmt.Printf("Feed created\n")
	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("updated_at: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("Url: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)
    return nil
}

func handlerGetFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
			fmt.Printf("Feed: %s\nURL: %s\nUser: %s\n\n", feed.Name, feed.Url, feed.CreatedBy)
	}
	return nil
}