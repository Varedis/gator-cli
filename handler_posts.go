package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/varedis/gator-cli/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("couldn't parse limit: %v", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %v", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %s\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("================================================")
	}

	return nil
}

func handlerSearch(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Usage: %s <search-term> [limit=2]\n", cmd.Name)
	}

	searchTerm := cmd.Args[0]
	limit := 2
	if len(cmd.Args) == 2 {
		if specifiedLimit, err := strconv.Atoi(cmd.Args[1]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("couldn't parse limit: %v", err)
		}
	}

	posts, err := s.db.FindPostsForUser(context.Background(), database.FindPostsForUserParams{
		UserID:  user.ID,
		Column2: sql.NullString{String: searchTerm, Valid: true},
		Limit:   int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %v", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %s\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("================================================")
	}

	return nil

	return nil
}
