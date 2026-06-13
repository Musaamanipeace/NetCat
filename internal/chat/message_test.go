package chat

import (
	"strings"
	"testing"
)

func TestFormatMessage(t *testing.T) {
	sender := "Alice"
	msg := "Hello world"
	result := FormatMessage(sender, msg)

	if !strings.Contains(result, "[Alice]: Hello world\n") {
		t.Errorf("FormatMessage layout mismatch, got: %q", result)
	}
}

func TestFormatSystemMessage(t *testing.T) {
	msg := "Bob has joined our chat."
	result := FormatSystemMessage(msg)

	if !strings.Contains(result, "[System]: Bob has joined our chat.\n") {
		t.Errorf("FormatSystemMessage layout mismatch, got: %q", result)
	}
}