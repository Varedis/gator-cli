package main

import (
	"fmt"
	"log"

	"github.com/varedis/gator-cli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("Read config %+v\n", cfg)

	if err = cfg.SetUser("Rob"); err != nil {
		log.Fatalf("couldn't set the current user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	fmt.Printf("Read config again %+v\n", cfg)
}
