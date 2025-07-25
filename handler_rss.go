package main

import (
	"context"
	"fmt"

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
