package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/varedis/gator-cli/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   name,
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("cannot create feed follow: %v", err)
	}

	fmt.Println("Feed Created:")
	fmt.Printf("  * ID: %s\n", feed.ID)
	fmt.Printf("  * Name: %s\n", feed.Name)
	fmt.Printf("  * URL: %s\n", feed.Url)
	fmt.Printf("  * Created at: %s\n", feed.CreatedAt)
	fmt.Printf("  * Updated at: %s\n", feed.CreatedAt)
	fmt.Printf("  * User ID: %s\n", feed.UserID)
	fmt.Println("Feed following successfully:")
	fmt.Printf("  * User: %s\n", feedFollow.UserID)
	fmt.Printf("  * Feed: %s\n", feedFollow.FeedName)
	fmt.Println("=============================")

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %v", err)
	}
	fmt.Println("Found feeds:")
	for _, feed := range feeds {
		fmt.Printf("  * Name: %s | URL: %s | User: %s\n", feed.Name, feed.Url, feed.User)
	}
	return nil
}
