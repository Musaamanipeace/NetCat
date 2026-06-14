// Package main is the entry point for the TCP Chat server application.
//
// This file is responsible for:
//   - Loading server configuration from command-line arguments
//   - Starting the TCP chat server in the background
//   - Keeping the program running until the user stops it
//
// Think of main.go as the front door of our chat building - everything starts here.
package main

// "fmt" - This package lets us write formatted text to the screen.
//   Example: fmt.Println("Hello") prints "Hello" followed by a newline.
// "os" - This package lets us interact with the operating system.
//   Example: os.Args gives us the command-line arguments.
//   Example: os.Exit(1) stops the program (we use it for errors).
// "net/internal/chat" - This is OUR custom package containing chat server logic.
//   It's a toolbox we built for handling rooms, clients, and messages.
// "net/pkg/config" - This is OUR custom package for loading server settings.
//   It handles command-line arguments and validates the port number.
import (
	"fmt"    // For printing messages to the console and to error output
	"os"     // For exiting the program and reading command-line arguments

	"net/internal/chat" // Our custom chat server logic
	"net/pkg/config"    // Our custom configuration loader
)

// main is where the program starts running. When you type "go run ." or "./net-cat",
// this function is the very first thing that executes. Think of this as the "starting line"
// of our program.
func main() {
	// Load() reads the command-line arguments and gives us back a configuration object.
	// For example, if you run "./net-cat 2525", Load() returns a Config with Port="2525".
	// Configuration contains the port number we should listen on.
	cfg, err := config.Load()

	// The "if err != nil" pattern in Go means "if there was an error".
	// err being "nil" means "no error" (nil is like "nothing" in Go).
	// If Load() returned an error, something was wrong with the arguments.
	if err != nil {
		// fmt.Fprintln writes to a specific destination (here, os.Stderr = error output).
		// Unlike fmt.Println which writes to standard output, this writes to the error channel.
		fmt.Fprintln(os.Stderr, err) // Print the error message to error output

		// os.Exit(1) stops the program immediately and returns error code 1.
		// Exit code 0 means success, any non-zero code means something went wrong.
		os.Exit(1) // Exit with error status
	}

	// go chat.Start(cfg.Port) STARTS a NEW goroutine (background task).
	// A goroutine is like a lightweight worker that runs independently.
	// While the goroutine runs in the background, we continue to the NEXT line.
	// This is why our server keeps running - it's working in the background!
	// Think of it as: "Start the chat server and let it run by itself."
	go chat.Start(cfg.Port)

	// select{} blocks the main goroutine forever without using CPU.
	// This keeps the program running indefinitely.
	// When the user presses Ctrl+C, the operating system will terminate the process.
	// We don't use "os/signal" because it's not in the allowed packages.
	select{}
}