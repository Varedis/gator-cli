package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/varedis/gator-cli/internal/database"
	"github.com/varedis/gator-cli/internal/rss"
)

const URL = "https://www.wagslane.dev/index.xml"

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("cannot parse the duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	// Find the next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get the next feed to fetch: %v", err)
	}

	// Retrieve the feed from the URL
	fetched, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("cannot fetch the feed: %v", err)
	}

	// Mark the feed as processed
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("cannot mark the feed as fetched: %v", err)
	}

	for _, item := range fetched.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
		parsedTime, err := time.Parse(time.RFC822, item.PubDate)
		if err != nil {
			return fmt.Errorf("couldn't parse the publish date: %v", err)
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: parsedTime, Valid: true},
			FeedID:      feed.ID,
		})
		if err != nil {
			return fmt.Errorf("error inserting post: %v", err)
		}
		fmt.Printf("Post Saved: %s\n", post.ID)
	}

	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(fetched.Channel.Item))

	return nil
}
