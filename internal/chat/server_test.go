// Package chat_test tests our chat server's core functionality.
//
// These tests make sure our server correctly:
//   - Rejects connections when the chat is full
//   - Allows connections when there's room
//   - Handles multiple clients at once
//
// We use "127.0.0.1" (localhost) to test without needing a real network.
package chat

// "bufio" - For reading lines of text from network connections.
// "net" - For simulating client connections to our server.
// "strings" - For checking if text contains certain phrases.
// "testing" - Go's built-in testing framework.
// "time" - For small delays (letting server process connections).
import (
	"bufio"   // Buffered I/O - reads text efficiently from connections
	"net"     // Network connections (net.Dial to simulate clients)
	"strings" // String tools (Contains to check for "has joined")
	"testing" // Go's testing framework
	"time"    // Time functions (time.Sleep for small delays)
)

// TestServer_MaxClientsLimit verifies that our server properly rejects
// connections when 10 clients are already chatting.
//
// This test is like: "Fill up the room, then try to squeeze in one more person!"
func TestServer_MaxClientsLimit(t *testing.T) {
	// Step 1: Get a free port for testing.
	// net.Listen("tcp", "127.0.0.1:0") asks the OS for an available port.
	// Port 0 means "give me any free port" - perfect for testing!
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		// t.Fatalf stops the test immediately with an error.
		t.Fatalf("failed to listen: %v", err)
	}
	// Extract the port number from the listener's address.
	// SplitHostPort separates "127.0.0.1:PORT" into host and port parts.
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	listener.Close() // We don't need this anymore - Start() creates its own

	// Step 2: Start our chat server in the background.
	// go makes Start() run concurrently while tests continue.
	go Start(port)

	// Step 3: Wait for the server to actually start listening.
	// TCP binding isn't instant, so we retry a few times.
	// This loop tries to connect until successful or 10 attempts pass.
	var initialConn net.Conn
	for i := 0; i < 10; i++ {
		initialConn, err = net.Dial("tcp", "127.0.0.1:"+port)
		if err == nil {
			initialConn.Close() // We just wanted to test if server is up
			break
		}
		// time.Sleep pauses for 10 milliseconds before retrying.
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("server failed to start on port %s: %v", port, err)
	}

	// Step 4: Connect maxClients (10) test clients to fill the room.
	// Each client will join and receive the "joined" system message.
	conns := make([]net.Conn, maxClients)
	scanners := make([]*bufio.Scanner, maxClients)

	for i := 0; i < maxClients; i++ {
		// Each iteration creates one client and connects them.
		conn, err := net.Dial("tcp", "127.0.0.1:"+port)
		if err != nil {
			t.Fatalf("failed to connect client %d: %v", i+1, err)
		}
		conns[i] = conn
		// t.Cleanup registers a function to run after the test ends.
		// This ensures we close all connections even if the test fails.
		t.Cleanup(func() { conn.Close() })

		// Each client sends a unique name: "A", "B", "C", etc.
		// rune(65+i) gives us the ASCII code for "A" (65), "B" (66), etc.
		_, err = conn.Write([]byte(string(rune(65+i)) + "\n"))
		if err != nil {
			t.Fatalf("failed handshake payload for client %d: %v", i+1, err)
		}

		// Create a scanner to read lines from this client's connection.
		// We use a larger buffer because welcome messages can be long.
		scanner := bufio.NewScanner(conn)
		scanner.Buffer(make([]byte, 0, 2*1024), 256*1024)
		scanners[i] = scanner

		// Wait for the "joined" confirmation message.
		// The server sends this when a client successfully joins.
		joinMsg := ""
		foundJoin := false
		for scanner.Scan() {
			line := scanner.Text()
			joinMsg += line + "\n"
			// Check if we got the join confirmation.
			if strings.Contains(line, "has joined our chat.") {
				foundJoin = true
				break
			}
		}
		if !foundJoin {
			t.Fatalf("client %d did not receive join confirmation", i+1)
		}
	}

	// Give the server a moment to finish registering all 10 clients.
	// Without this delay, the next connection might sneak in before we're full.
	time.Sleep(200 * time.Millisecond)

	// Step 5: Try to connect one more client (should be rejected).
	var overflowConn net.Conn
	var response string

	// Try a few times since network operations can be unreliable.
	for attempts := 0; attempts < 10; attempts++ {
		overflowConn, err = net.Dial("tcp", "127.0.0.1:"+port)
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Read the response from the rejected connection.
		// When rejected, we should see "Chat is full" message.
		scanner := bufio.NewScanner(overflowConn)
		if scanner.Scan() {
			response = strings.TrimSpace(scanner.Text())
			overflowConn.Close()
			break
		}
		overflowConn.Close()
		time.Sleep(10 * time.Millisecond)
	}

	// Step 6: Verify the rejection message is correct.
	expected := "Chat is full. Try again later."
	if response != expected {
		t.Errorf("expected rejection message %q, got %q", expected, response)
	}
}