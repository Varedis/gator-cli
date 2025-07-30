package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/varedis/gator-cli/internal/database"
	"github.com/varedis/gator-cli/internal/rss"
)

const URL = "https://www.wagslane.dev/index.xml"

func handlerRSS(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), URL)
	if err != nil {
		return fmt.Errorf("couldn't retrieve rss feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()

	currentUser := s.cfg.CurrentUserName
	user, err := s.db.GetUser(ctx, currentUser)
	if err != nil {
		return fmt.Errorf("cannot get current user: %v", err)
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   name,
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %v", err)
	}

	fmt.Println("Feed Created:")
	fmt.Printf("  * ID: %s\n", feed.ID)
	fmt.Printf("  * Name: %s\n", feed.Name)
	fmt.Printf("  * URL: %s\n", feed.Url)
	fmt.Printf("  * Created at: %s\n", feed.CreatedAt)
	fmt.Printf("  * Updated at: %s\n", feed.CreatedAt)
	fmt.Printf("  * User ID: %s\n", feed.UserID)

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

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	ctx := context.Background()

	// Get Feed
	feedID, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("cannot get feed: %v", err)
	}

	currentUser := s.cfg.CurrentUserName
	user, err := s.db.GetUser(ctx, currentUser)
	if err != nil {
		return fmt.Errorf("cannot get current user: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feedID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %v", err)
	}

	fmt.Println("Followed Feed:")
	fmt.Printf("  * Name: %s\n", feedFollow.FeedName)
	fmt.Printf("  * User: %s\n", feedFollow.UserName)

	return nil
}
