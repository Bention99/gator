package main

import (
	"errors"
	"fmt"
	"time"
	"database/sql"
	"context"
	"strings"
	"log"
	"strconv"
	"github.com/google/uuid"
	"github.com/Bention99/gator/internal/config"
	"github.com/Bention99/gator/internal/database"
	"github.com/Bention99/gator/internal/rss"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err != nil {
			return errors.New("No Account for that username.")
		}
		err = handler(s, cmd, user)
		return err
	}
}

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
	if len(cmd.args) == 0 {
        return errors.New("agg expects a single argument: time_between_reqs, like: 1ms, 1m, 1h")
    }

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		return err
	}
	feed, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	for i := range feed.Channel.Item {
		now := sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

		var publishedAt sql.NullTime
		if feed.Channel.Item[i].PubDate != "" {
			t, err := time.Parse(time.RFC1123, feed.Channel.Item[i].PubDate)
			if err != nil {
				t, err = time.Parse(time.RFC1123Z, feed.Channel.Item[i].PubDate)
			}
			if err == nil {
				publishedAt = sql.NullTime{
					Time:  t,
					Valid: true,
				}
			}
		}

		cpp := database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title: feed.Channel.Item[i].Title,
			Url: sql.NullString{
				String: feed.Channel.Item[i].Link,
				Valid:  feed.Channel.Item[i].Link != "",
			},
			Description: sql.NullString{
				String: feed.Channel.Item[i].Description,
				Valid:  feed.Channel.Item[i].Description != "",
			},
			PublishedAt: publishedAt,
			FeedID: nextFeed.ID,
		}

		_, err := s.db.CreatePost(ctx, cpp)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
        		continue
			} else {
				log.Printf("Error creating post: %v", err)
				return err
			}
		}
		fmt.Printf("Post saved: %v", feed.Channel.Item[i].Title)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command, currentUser database.User) error {
	ctx := context.Background()

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

	cffp := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID: feed.UserID,
		FeedID: feed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, cffp)
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

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	ctx := context.Background()

	feed, err := s.db.GetFeed(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	now := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	cffp := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID: currentUser.ID,
		FeedID: feed.ID,
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, cffp)
	if err != nil {
		return err
	}
	fmt.Printf("User: %s\nis now following Feed: %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, currentUser database.User) error {
	ctx := context.Background()

	feedRows, err := s.db.GetFeedFollowsForUser(ctx, currentUser.Name)
	if err != nil {
		return err
	}

	if len(feedRows) == 0 {
		fmt.Println("You are currently not following any feeds.")
	}

	for _, row := range feedRows {
		fmt.Printf("Feed: %s\nUser: %s\n", row.FeedName, row.UserName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, currentUser database.User) error {
	ctx := context.Background()

	ufp := database.UnfollowFeedParams {
		Url: cmd.args[0],
		Name: currentUser.Name,
	}

	err := s.db.UnfollowFeed(ctx, ufp)
	if err != nil {
		return err
	}
	fmt.Printf("You unfollowed: %v\n", cmd.args[0])
	return nil
}

func handlerBrowse(s *state, cmd command) error {
	ctx := context.Background()
	var limit int32 = 2
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = int32(parsedLimit)
	}
	posts, err := s.db.GetPostsForUser(ctx, limit)
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("%v\n", post.Title)
		fmt.Printf("%v\n", post.Url)
		fmt.Printf("%v\n", post.Description)
		fmt.Printf("%v\n", post.PublishedAt)
	}
	return nil
}