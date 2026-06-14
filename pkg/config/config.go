// Package config handles loading and validating server settings from command-line arguments.
//
// Think of this package as the "settings checker" for our chat server.
// It looks at what the user typed when starting the program and figures out
// things like: which port should the server listen on?
//
// For example:
//   ./net-cat           -> uses default port 8989
//   ./net-cat 2525      -> uses port 2525
//   ./net-cat 2525 x    -> ERROR! Too many arguments
package config

// "fmt" - This package lets us create error messages with formatting.
//   Example: fmt.Errorf("error: %v", value) creates a formatted error.
// "os" - This package lets us read command-line arguments the user typed.
//   Example: os.Args is a list like ["./net-cat", "2525"]
// "strconv" - This package converts between strings and numbers.
//   Example: strconv.Atoi("2525") gives us the number 2525
import (
	"fmt"     // For creating error messages
	"os"      // For reading command-line arguments (os.Args)
	"strconv" // For converting string to integer (Atoi = Ascii to Integer)
)

// Config is a container (like a form) that holds all our server settings.
// Think of it as a "settings sheet" with all the important information we need.
// In Go, you define a "struct" (short for "structure") to group related data together.
//
// Fields explained:
//   - Port: Where the server listens (like a door number). Default is 8989.
//   - MaxClients: Maximum number of people who can join (default is 10).
//   - LogFile: Where we store the chat history (empty means no logging).
type Config struct {
	Port      string // Door number for the server (like house number)
	MaxClients int   // Maximum people in the chat room at once
	LogFile   string // File path to save chat logs (empty = no file logging)
}

// Load reads the command-line arguments and returns a valid configuration.
// It validates the port number and returns an error if something is wrong.
//
// Command-line arguments explained:
//   - os.Args[0] is the program name (like "./net-cat")
//   - os.Args[1] is the first argument after the program name (like "2525")
//   - Example: Running "./net-cat 2525" gives os.Args = ["./net-cat", "2525"]
func Load() (*Config, error) {
	// Create a Config with sensible defaults.
	// If user doesn't specify anything, these values are used.
	// This is like having a "fallback plan" if the user doesn't give us specific instructions.
	cfg := &Config{
		Port:      "8989", // Default door number if user doesn't specify one
		MaxClients: 10,    // Maximum 10 people can chat at once (per room)
		LogFile:   "",      // No log file by default (empty string = disabled)
	}

	// Check how many arguments the user typed.
	// len(os.Args) counts the total number.
	// If user types "./net-cat 2525 abc", len(os.Args) = 3 (program + 2 args)
	if len(os.Args) > 2 {
		// We got too many arguments! According to the rules, we should show:
		// [USAGE]: ./TCPChat $port
		return nil, fmt.Errorf("[USAGE]: ./TCPChat $port")
	}

	// If there's exactly 1 argument, user is probably specifying a port
	// (the original behavior before we added more features)
	if len(os.Args) == 2 {
		// strconv.Atoi tries to convert the text into a number.
		// For example, "2525" becomes the number 2525.
		// If it's not a number (like "abc"), it returns an error.
		portNum, err := strconv.Atoi(os.Args[1])
		if err != nil || portNum < 1 || portNum > 65535 {
			// Three cases we're checking:
			// 1. err != nil: User typed something that's not a number
			// 2. portNum < 1: Port 0 or negative doesn't make sense
			// 3. portNum > 65535: Port numbers only go up to 65535 (internet rule)
			return nil, fmt.Errorf("[USAGE]: ./TCPChat $port\nPort must be a number between 1 and 65535")
		}
		// All good! Use the port number user specified.
		cfg.Port = os.Args[1]
	}

	// Return the configuration. The user can now use cfg.Port in their code.
	return cfg, nil
}