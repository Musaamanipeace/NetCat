// Package chat manages individual chat clients (people connecting to our server).
//
// A "client" is someone connected to the chat via `nc` or similar tools.
// Each client has their own connection, name, and message mailbox (channel).
//
// Think of a client as one person in a room with a walkie-talkie.
// They can speak (send messages) and listen (receive messages from others).
package chat

// "bufio" - For reading lines of text efficiently from connections.
// "errors" - For creating simple error values without formatting.
// "net" - For network connections (the "pipe" between server and client).
// "strings" - For trimming whitespace from names and messages.
import (
	"bufio"  // Buffered I/O - reads/writes text efficiently
	"errors" // Simple error creation (errors.New)
	"net"    // Network connections (TCP sockets)
	"strings" // String tools (TrimSpace removes spaces from beginning/end)
)

// Client represents one person who has connected to our chat server.
// Each time someone connects with `nc localhost 8989`, a Client object is created.
//
// Think of this like a "profile" for each person in the chat - it stores their
// connection details, their name, and how to send them messages.
//
// Fields explained:
//   - conn: The "phone line" to this specific client (how we talk to them)
//   - name: The name they chose when joining (like "Alice" or "Bob")
//   - server: Pointer back to the chat server (for broadcasting messages)
//   - ch: A "mailbox" channel where we put messages to send them
//   - scanner: A tool that reads their messages line by line
//   - writer: A tool that writes messages back to them
type Client struct {
	conn    net.Conn        // The network "phone line" to this client
	name    string          // The client's chosen name in the chat
	server  *Server        // Pointer to the server running this chat
	ch      chan string     // Their personal "inbox" - messages come here
	scanner *bufio.Scanner // Reads what they type, one line at a time
	writer  *bufio.Writer   // Writes messages back to their screen
}

// NewClient creates a new Client and prepares them for chatting.
// This is called every time a new person connects to our server.
//
// Think of this as: "Welcome! Here's your chat profile and tools."
//
// Parameters:
//   - conn: Their network connection (the phone line)
//   - server: The chat server they're joining
//
// Returns:
//   - A pointer to the newly created Client
func NewClient(conn net.Conn, server *Server) *Client {
	// Return a Client struct with all fields filled in.
	// &Client{...} creates a pointer to this new Client.
	// net.Conn is the connection we pass in.
	// make(chan string, 100) creates a buffered channel with room for 100 messages.
	// Buffered means: if 5 messages come in quickly, they're stored until we read them.
	// bufio.NewScanner reads text one line at a time (like reading a book line by line).
	// bufio.NewWriter buffers text before sending (more efficient than sending character by character).
	return &Client{
		conn:    conn,               // Save the connection
		server:  server,             // Save the server reference
		ch:      make(chan string, 100), // Create their message inbox
		scanner: bufio.NewScanner(conn),  // Set up a scanner to read their messages
		writer:  bufio.NewWriter(conn),   // Set up a writer to send messages to them
	}
}

// Read is a goroutine that continuously listens for messages from this client.
// A goroutine is a lightweight background task that runs independently.
// This means: while we're reading THIS client's messages, other clients can still chat!
//
// Think of this as: "Sit and listen to what this person is saying."
func (c *Client) Read() {
	// for c.scanner.Scan() loops until the client disconnects.
	// scanner.Scan() reads one line of text and returns true if successful.
	// When the client closes their terminal or types Ctrl+C, it returns false.
	for c.scanner.Scan() {
		// .Text() gets the line they typed (as a string).
		msg := strings.TrimSpace(c.scanner.Text())

		// strings.TrimSpace removes spaces/tabs/newlines from beginning and end.
		// Example: "  hello  \n" becomes "hello"

		// Empty messages (just hitting Enter) are ignored.
		// The server doesn't broadcast these to avoid spam.
		if msg == "" {
			continue // Skip to the next message
		}

		// Check if this is a special command (like /name).
		// Commands start with "/" and do special actions instead of normal chatting.
		if strings.HasPrefix(msg, "/") {
			c.handleCommand(msg) // Process the command (like changing name)
			continue             // Don't broadcast commands as regular messages
		}

		// Broadcast the message to all other clients in the chat.
		// c.server.broadcast does the heavy lifting of sending it to everyone else.
		c.server.broadcast(msg, c)
	}

	// If we reach this line, the client disconnected (closed their terminal).
	// We tell the server to clean up and notify others.
	c.server.removeClient(c)
}

// handleCommand processes special commands the client typed (starting with /).
// These are actions like changing names, not regular chat.
//
// Supported commands:
//   - /name <newname>: Changes the client's displayed name
func (c *Client) handleCommand(cmd string) {
	// Split the command into parts. "strings.Fields splits by whitespace.
	// Example: "/name Alice" becomes ["name", "Alice"]
	parts := strings.Fields(cmd)

	// We need at least 2 parts: the command name and its argument.
	if len(parts) < 2 {
		c.ch <- FormatSystemMessage("Usage: /name <newname>")
		return
	}

	// The first part after "/" is the command name (without the slash).
	command := parts[0][1:]   // Remove the leading "/" from "/name" to get "name"
	arg := parts[1]            // The argument (like "Alice")

	switch command {
	case "name":
		// Change the client's name.
		c.changeName(arg)
	default:
		// Unknown command - tell them what's available.
		c.ch <- FormatSystemMessage("Unknown command. Try /name <newname>")
	}
}

// changeName lets a client pick a new name during the chat.
// This is a bonus feature - normally names are permanent after joining.
func (c *Client) changeName(newName string) {
	// Don't allow empty names.
	if newName == "" {
		c.ch <- FormatSystemMessage("Name cannot be empty!")
		return
	}

	// Ask the server to change our name (server handles the actual update).
	oldName := c.server.renameClient(c, newName)

	// If the rename succeeded, oldName contains their previous name.
	if oldName != "" {
		// Tell everyone in the chat about the name change.
		msg := FormatSystemMessage(oldName + " is now known as " + newName)
		c.server.messages <- msg
	} else {
		// Name was already taken.
		c.ch <- FormatSystemMessage("Name already taken. Try another name.")
	}
}

// Write is a goroutine that continuously sends messages to this client.
// This runs in the background, feeding messages from their inbox (ch) to their screen.
//
// Think of this as: "Take messages from their mailbox and deliver them."
func (c *Client) Write() {
	// for msg := range c.ch waits for new messages to arrive in their inbox.
	// When the server sends a message to this client, it appears here.
	for msg := range c.ch {
		// WriteString adds text to the buffer (doesn't send yet).
		// Flush actually sends all buffered text to the client's screen.
		c.writer.WriteString(msg) // Add to buffer
		c.writer.Flush()          // Send the buffer contents
	}
	// When c.ch is closed (in Close()), this loop ends automatically.
}

// Close cleans up when a client leaves or gets disconnected.
// This is like "hanging up the phone" - we close their connection.
func (c *Client) Close() {
	// Close the network connection - this stops any future communication.
	c.conn.Close()

	// Close the message channel - this tells the Write() goroutine to stop.
	// Closing a channel is like "locking the mailbox" - no more messages can be delivered.
	close(c.ch)
}

// welcomeMessage sends the welcome logo and asks for the client's name.
// This is the first thing a new client sees when they connect.
//
// Think of this as the "hello and registration" step.
//
// Returns:
//   - error if something went wrong (like connection closed before they could send name)
func (c *Client) welcomeMessage() error {
	// First, send the ASCII art logo (a cute penguin saying "Welcome to TCP-Chat!").
	// GetWelcomeLogo() returns the penguin art + the name prompt.
	c.writer.WriteString(GetWelcomeLogo()) // GetWelcomeLogo() returns the penguin art + prompt
	c.writer.Flush()                     // Actually send it to the client

	// Wait for the client to type their name and press Enter.
	// scanner.Scan() returns true if they sent a line of text.
	if c.scanner.Scan() {
		// Get what they typed and remove any extra spaces.
		name := strings.TrimSpace(c.scanner.Text())

		// Names can't be empty (you must have a name to chat).
		if name == "" {
			c.writer.WriteString("Name cannot be empty. Disconnecting.\n")
			c.writer.Flush()
			// Return an error to signal that setup failed.
			return errors.New("empty name") // errors.New creates a simple error
		}
		// Save their chosen name in their profile.
		c.name = name
	} else {
		// scanner.Scan() returned false - maybe they disconnected immediately?
		return errors.New("failed to read name")
	}
	return nil // No error means success
}