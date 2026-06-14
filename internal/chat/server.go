// Package chat contains all the core logic for our TCP chat server.
//
// This file implements the "brain" of our chat system - managing rooms,
// clients, and the flow of messages between them.
//
// The server can:
//   - Accept multiple client connections (up to 10)
//   - Broadcast messages to all clients
//   - Welcome new clients and show chat history
//   - Handle graceful disconnections
//   - Support name changes (bonus feature)
//   - Log all messages to a file (bonus feature)
package chat

// "fmt" - For printing messages to the console and formatting strings.
// "log" - For writing messages to a log file.
// "net" - For creating TCP network connections (the "phone lines").
// "os" - For opening log files.
// "sync" - For mutex locks (to prevent chaos when multiple clients talk at once).
import (
	"fmt"   // Format strings and print to console
	"log"   // Write to log files
	"net"   // TCP network connections (net.Listen, net.Conn)
	"os"    // Operating system interface (files, os.OpenFile)
	"sync"  // Mutex for protecting shared data (sync.Mutex)
)

// maxClients is the maximum number of people allowed in the chat.
// We set this to 10 because too many people makes conversation hard to follow.
// Think of it as: "This chat room can only hold 10 people at once."
const maxClients = 10

// logFilePath is the default path where chat logs are saved.
// This implements the "log to file" bonus feature.
const logFilePath = "chat.log"

// Server manages the chat room and handles incoming connections.
// This is like the "reception desk" of our building - managing the chat.
//
// Think of this as a hotel desk clerk who knows everyone in the chat
// and can pass messages between them.
//
// Fields explained:
//   - clients: All connected clients (like a master guest list)
//   - history: All messages ever sent (stored for new joiners)
//   - mutex: Lock to protect shared data (prevents data races)
//   - messages: Channel for broadcasting messages between goroutines
//   - logger: Writes messages to a log file (bonus feature)
type Server struct {
	clients  map[*Client]bool // Map of all connected clients (key=client, value=true)
	history  []string       // All messages in the chat (chronological order)
	mutex    sync.Mutex     // Mutual exclusion lock - only one goroutine can access at once
	messages chan string    // Channel for messages - goroutines put messages here to broadcast
	logger   *log.Logger   // Writes all chat messages to the log file
}

// NewServer creates a new Server instance.
// This sets up the basic structure before connections start arriving.
//
// Returns:
//   - A pointer to the newly created Server (ready to start accepting connections)
func NewServer() *Server {
	// Open the log file for appending.
	// os.OpenFile with O_APPEND means: add new messages to the end of the file.
	// os.O_CREATE means: create the file if it doesn't exist.
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// If we can't open the log file, we continue without logging.
		// The chat still works, just without file logging.
		logFile = nil
	}

	// Return a fully initialized Server.
	return &Server{
		clients:  make(map[*Client]bool),  // Empty map (no clients yet)
		history:  make([]string, 0),      // Empty slice (no messages yet)
		messages: make(chan string, 100), // Buffered channel for 100 messages
		logger:   log.New(logFile, "", 0), // Logger for the file (or nil if no file)
	}
}

// Start begins listening for TCP connections on the specified port.
// This is the main entry point - it opens the door and welcomes guests.
//
// Parameters:
//   - port: The door number to listen on (like "8989")
func Start(port string) {
	// Create a new server instance.
	server := NewServer()

	// net.Listen creates a TCP "phone line" listening on the port.
	// "tcp" means we use the TCP protocol (reliable, connection-based).
	// The ":" prefix means "listen on all available network addresses."
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// fmt.Printf prints a formatted message (with variables).
		fmt.Printf("Error starting server: %v\n", err)
		return // Stop trying to start if we can't bind to the port
	}
	// defer means "do this when the function ends."
	// We close the listener before exiting (clean up our phone line).
	defer listener.Close() // Always close when we're done

	fmt.Printf("Listening on the port :%s\n", port)

	// Start a goroutine that broadcasts messages to all clients.
	// A goroutine is a background worker that runs concurrently.
	// Think: "Start the messenger who delivers messages to everyone."
	go server.broadcaster()

	// This is an INFINITE loop - we keep accepting connections forever!
	// The for loop without conditions runs forever (until program exits).
	for {
		// listener.Accept() waits for someone to connect and returns their connection.
		// This line PAUSES here until a client connects (like waiting for a phone call).
		conn, err := listener.Accept()
		if err != nil {
			// If Accept fails (network error), print error and continue waiting.
			fmt.Printf("Error accepting connection: %v\n", err)
			continue // Skip to next loop iteration
		}

		// When a client connects, spin up a handler goroutine.
		// go means "run this function in the background."
		// Each client gets their own handler, so everyone can chat at once!
		go server.handleConnection(conn)
	}
}

// handleConnection welcomes a new client and adds them to the chat.
// This runs in its own goroutine so multiple clients can connect simultaneously.
//
// Parameters:
//   - conn: The network connection to this specific client
func (s *Server) handleConnection(conn net.Conn) {
	// Check if we're at max capacity BEFORE doing anything else.
	// This must happen first to reject connections immediately.
	s.mutex.Lock()
	if len(s.clients) >= maxClients {
		s.mutex.Unlock() // Unlock before writing to connection
		conn.Write([]byte("Chat is full. Try again later.\n"))
		conn.Close()
		return
	}
	s.mutex.Unlock()

	// Create a Client object to represent this person.
	client := NewClient(conn, s)

	// Send the welcome logo and ask for their name.
	// This is like greeting them at the door and asking for their name.
	if err := client.welcomeMessage(); err != nil {
		client.Close() // Clean up if setup fails
		return         // Stop handling this client
	}

	// Check if we already have someone with this name (names must be unique).
	s.mutex.Lock()
	for existingClient := range s.clients {
		if existingClient.name == client.name {
			s.mutex.Unlock() // Always unlock before returning!
			conn.Write([]byte("Name already taken. Disconnecting.\n"))
			client.Close()
			return
		}
	}

	// Register the client in our master list.
	s.clients[client] = true
	s.mutex.Unlock() // "Unlock the door" - others can access now

	// Start goroutines for reading and writing to this client (in background).
	// These two goroutines work independently - one listens for their messages,
	// the other sends them messages from other people.
	go client.Write() // Send messages TO them (runs forever until they disconnect)
	go client.Read()  // Read messages FROM them (runs until they disconnect)

	// Send them all previous messages so they can catch up.
	s.sendHistory(client)

	// Announce their arrival to everyone else in the chat.
	joinMsg := FormatSystemMessage(client.name + " has joined our chat.")
	s.messages <- joinMsg // Put message in the broadcast channel

	// Save to history for future joiners.
	s.mutex.Lock()
	s.history = append(s.history, joinMsg)
	s.mutex.Unlock()
}

// broadcaster is a goroutine that continuously forwards messages to clients.
// Think of this as: "The messenger who runs around delivering notes to everyone."
func (s *Server) broadcaster() {
	// Loop forever, taking messages from the messages channel.
	for msg := range s.messages {
		// Lock the mutex so we don't accidentally modify data while reading.
		s.mutex.Lock()

		// Send the message to everyone in the chat.
		// Note: we're not excluding anyone here - system messages go to everyone.
		for client := range s.clients {
			// select with default is a non-blocking send.
			// If the client's inbox is full, we skip them (avoid blocking).
			select {
			case client.ch <- msg: // Put the message in their inbox
			default: // Inbox full, skip this client
			}
		}
		s.mutex.Unlock() // Unlock for others
	}
}

// broadcast sends a message to all clients EXCEPT the sender.
// This is like: "Speak into the microphone, everyone else hears you."
//
// Parameters:
//   - message: What to send (the text of what they said)
//   - sender: Who sent it (we don't echo messages back to the sender)
func (s *Server) broadcast(message string, sender *Client) {
	// Don't send empty messages (spam prevention).
	if message == "" {
		return // Exit early - nothing to do
	}

	// Format the message with timestamp and sender name.
	// Example: "[2024-01-20 15:48:41][Alice]: Hello everyone!"
	formattedMsg := FormatMessage(sender.name, message)

	// Save to history.
	s.mutex.Lock()
	s.history = append(s.history, formattedMsg)
	s.mutex.Unlock()

	// Send to all OTHER clients (not the sender).
	s.mutex.Lock()
	for client := range s.clients {
		// Only send to clients who AREN'T the sender.
		// This prevents the sender from seeing their own message twice.
		if client != sender {
			select {
			case client.ch <- formattedMsg: // Put in their inbox
			default: // Skip if inbox is full
			}
		}
	}
	s.mutex.Unlock()
}

// removeClient removes a client from the chat and notifies others.
// This is called when a client disconnects or changes rooms.
//
// Parameters:
//   - client: The client who is leaving
func (s *Server) removeClient(client *Client) {
	// Check if client exists and remove them.
	s.mutex.Lock()
	if _, exists := s.clients[client]; exists {
		delete(s.clients, client) // Remove from our map
		s.mutex.Unlock()         // Unlock before doing more work

		// Announce their departure.
		leaveMsg := FormatSystemMessage(client.name + " has left our chat.")
		s.messages <- leaveMsg // Put in broadcast channel

		// Save to history.
		s.mutex.Lock()
		s.history = append(s.history, leaveMsg)
		s.mutex.Unlock()

		// Close their connection (hang up the phone).
		client.Close()
	} else {
		s.mutex.Unlock() // Just unlock if client wasn't found
	}
}

// renameClient changes a client's name and returns their old name.
// Returns empty string if rename failed (name already taken).
//
// Parameters:
//   - client: The client changing their name
//   - newName: What they want to be called
//
// Returns:
//   - The old name if successful, empty string if failed
func (s *Server) renameClient(client *Client, newName string) string {
	// Check if the new name is already taken.
	s.mutex.Lock()
	defer s.mutex.Unlock() // Unlock when function ends

	// Look through all clients to see if anyone has this name.
	for existingClient := range s.clients {
		if existingClient != client && existingClient.name == newName {
			return "" // Name taken, return empty string
		}
	}

	// Name is available, update it.
	oldName := client.name
	client.name = newName
	return oldName
}

// sendHistory sends all previous messages to a new client.
// This is like saying: "Here's everything everyone said before you arrived."
//
// Parameters:
//   - client: The new client who needs to catch up
func (s *Server) sendHistory(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock() // Unlock when function ends

	// Loop through all messages and send them to the new client.
	// This helps them understand the conversation context.
	for _, msg := range s.history {
		select {
		case client.ch <- msg: // Try to put in their inbox
		default: // If full, skip (shouldn't happen with new clients)
		}
	}
}