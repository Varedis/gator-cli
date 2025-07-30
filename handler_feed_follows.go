package main

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/varedis/gator-cli/internal/database"
)

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

func handlerFollowing(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUserName
	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("cannot get current user: %v", err)
	}

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("cannot get feed follows for user: %v", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user")
		return nil
	}

	fmt.Println("Following feeds:")
	for _, ff := range feedFollows {
		fmt.Printf("  * Name: %s\n", ff.FeedName)
	}

	return nil
}
