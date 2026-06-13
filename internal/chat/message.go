package chat

import (
	"fmt"
	"time"
)

// FormatMessage structures regular chat messages with timestamps.
func FormatMessage(sender string, msg string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]: %s\n", timestamp, sender, msg)
}

// FormatSystemMessage structures administrative/server broadcast items.
func FormatSystemMessage(msg string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][System]: %s\n", timestamp, msg)
}