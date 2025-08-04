package main

import (
	"context"
	"fmt"
	"time"

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
	}

	fmt.Printf("Feed %s collected, %v posts found\n", feed.Name, len(fetched.Channel.Item))

	return nil
}
