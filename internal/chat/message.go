// Package chat contains all the core logic for our TCP chat server.
//
// This package handles:
//   - Managing chat rooms where clients connect and chat
//   - Accepting new client connections
//   - Broadcasting messages between clients
//   - Handling client disconnections gracefully
//
// Think of this package as the "brain" of our chat application.
// It knows how to run chat rooms, add/remove clients, and pass messages around.
package chat

// "fmt" - For formatting text (like "[2024-01-20][Name]: Hello, world!")
// "time" - For getting the current date and time (when messages were sent)
import (
	"fmt"  // Combines text and variables into formatted strings
	"time" // Gets current time and formats it nicely
)

// FormatMessage structures a regular chat message with time and sender name.
// This function takes the sender's name and what they said, and wraps it
// in the required format: [2024-01-01 12:00:00][Alice]: Hello!
//
// Parameters explained:
//   - sender: The name of the person who sent the message (like "Alice")
//   - msg: What they typed (like "Hello everyone!")
//
// Returns:
//   - A formatted string ready to display in the chat
func FormatMessage(sender string, msg string) string {
	// time.Now() gives us the exact moment when this function runs.
	// Format("2006-01-02 15:04:05") converts the time into a string.
	// The numbers 2006-01-02 are special reference values Go uses for formatting.
	// Example result: "2024-01-20 15:48:41"
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// fmt.Sprintf combines text and variables into one string.
	// The %s placeholders get replaced with the actual values.
	// Example: If sender="Alice" and msg="Hello", result is:
	// "[2024-01-20 15:48:41][Alice]: Hello\n"
	// The \n adds a newline at the end (so each message appears on its own line).
	return fmt.Sprintf("[%s][%s]: %s\n", timestamp, sender, msg)
}

// FormatSystemMessage structures administrative messages from the server.
// These are messages that come from the "System" (the server itself), not from users.
// Example: "[2024-01-20 15:48:41][System]: Bob has joined our chat.\n"
//
// Parameters explained:
//   - msg: What happened (like "Bob has joined our chat.")
//
// Returns:
//   - A formatted system message ready to display
func FormatSystemMessage(msg string) string {
	// Same timestamp logic as FormatMessage.
	// Every message (system or user) gets the same timestamp format.
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// Format the message with "System" as the sender.
	// This makes it clear to users that this message came from the server.
	// Example: "has joined" becomes "[2024-01-20 15:48:41][System]: has joined our chat.\n"
	return fmt.Sprintf("[%s][System]: %s\n", timestamp, msg)
}