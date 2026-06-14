// Package chat_test tests our message formatting functions.
//
// These tests make sure messages look exactly right when they're displayed.
// We check that timestamps, names, and formatting all work correctly.
package chat

// "strings" - For checking if text contains certain pieces (HasPrefix, Contains, etc.)
// "testing" - For running test functions and asserting expectations
import (
	"strings" // String manipulation functions (Contains checks if text has a substring)
	"testing" // Go's built-in testing package
)

// TestFormatMessage checks that our regular chat messages are formatted correctly.
// When a user sends "Hello world", we expect to see something like:
// "[2024-01-20 15:48:41][Alice]: Hello world\n"
func TestFormatMessage(t *testing.T) {
	// Define who is sending the message and what they're saying.
	// Think of this like writing a note to a friend.
	sender := "Alice"    // The name of the person sending the message
	msg := "Hello world" // What they typed

	// Call our formatting function to get the nicely formatted message.
	result := FormatMessage(sender, msg)

	// Check that the result contains what we expect.
	// strings.Contains looks for one piece of text inside another.
	// We don't check the exact timestamp because that changes every time!
	expectedPart := "[Alice]: Hello world\n"

	// If the result doesn't contain our expected text, the test fails.
	if !strings.Contains(result, expectedPart) {
		// t.Errorf records the failure but lets the test continue.
		// The %q prints the value in quotes, making it easier to see.
		t.Errorf("FormatMessage layout mismatch, got: %q", result)
	}
}

// TestFormatSystemMessage checks that system messages (like join notifications) are formatted correctly.
// When someone joins, we expect: "[2024-01-20 15:48:41][System]: Bob has joined our chat.\n"
func TestFormatSystemMessage(t *testing.T) {
	// The message describes what happened in the chat.
	// In this case, Bob joining is a "system event" (not a chat message).
	msg := "Bob has joined our chat."

	// Format the system message.
	result := FormatSystemMessage(msg)

	// Check that we got the "System:" label and the message.
	// System messages are different from regular messages - they come from the server!
	expectedPart := "[System]: Bob has joined our chat.\n"

	if !strings.Contains(result, expectedPart) {
		t.Errorf("FormatSystemMessage layout mismatch, got: %q", result)
	}
}