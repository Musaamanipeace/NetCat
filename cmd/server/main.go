package main

import (
	"fmt"
	"os"

	"net/internal/chat"
	"net/pkg/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	chat.Start(cfg.Port)
}