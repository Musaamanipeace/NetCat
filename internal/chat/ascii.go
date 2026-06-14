// Package chat contains all the core logic for our TCP chat server.
//
// This file defines the welcome logo that new clients see when they connect.
// It's the ASCII art penguin that greets everyone at the door.
package chat

// GetWelcomeLogo returns the structured Penguin welcome banner for the TCP chat.
// When a client connects, this is the first thing they see in their terminal.
// It's a friendly way to welcome them and ask for their name.
//
// The logo is made of text characters arranged to look like a penguin:
//   - Lines starting with spaces create the shape
//   - Underscores, vertical bars, and other symbols form the penguin outline
//   - The "\n" at the end of each line creates a new line (like pressing Enter)
//   - The backslash "\" is escaped as "\\" because "\" has special meaning in Go strings
//
// Returns:
//   - A string containing the full welcome message including the name prompt
func GetWelcomeLogo() string {
	// Start with the welcome text.
	// The "\n" adds a newline character (moves cursor to next line).
	logo := "Welcome to TCP-Chat!\n"

	// Add each line of the penguin art.
	// Each "logo +=" appends more text to what we already have.
	// These lines were carefully designed to form a penguin shape.
	// Think of this like assembling a puzzle - each piece adds to the picture.
	logo += "         _nnnn_\n"      // Penguin's head outline
	logo += "        dGGGGMMb\n"    // Penguin's face
	logo += "       @p~qp~~qMb\n"   // Penguin's eyes and beak
	logo += "       M|@||@) M|\n"   // Penguin's body
	logo += "       @,----.JM|\n"   // More of penguin's body
	logo += "      JS^\\__/  qKL\n"  // Penguin's belly (note: \\ escapes to single \ in output)
	logo += "     dZP        qKRb\n" // Penguin's feet area
	logo += "    dZP          qKKb\n"
	logo += "   fZP            SMMb\n"
	logo += "   HZM            MMMM\n"
	logo += "   FqM            MMMM\n"
	logo += " __| \".        |\\dS\"qML\n" // Bottom of penguin
	logo += " |    `.       | `' \\Zq\n"
	logo += "_)      \\.___,|     .'\n"
	logo += "\\____   )MMMMMP|   .'\n"
	logo += "     `-'       `--'\n"

	// Finally, ask for their name. This prompts the user to type their chat name.
	logo += "[ENTER YOUR NAME]: "

	// Return the complete welcome message.
	return logo
}