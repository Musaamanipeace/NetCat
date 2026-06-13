package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds server configuration parameters.
type Config struct {
	Port string
}

// Load validates command-line arguments and returns a valid configuration.
func Load() (*Config, error) {
	port := "8989" // default port

	if len(os.Args) > 2 {
		return nil, fmt.Errorf("[USAGE]: ./TCPChat $port")
	}

	if len(os.Args) == 2 {
		portNum, err := strconv.Atoi(os.Args[1])
		if err != nil || portNum < 1 || portNum > 65535 {
			return nil, fmt.Errorf("[USAGE]: ./TCPChat $port\nPort must be a number between 1 and 65535")
		}
		port = os.Args[1]
	}

	return &Config{Port: port}, nil
}