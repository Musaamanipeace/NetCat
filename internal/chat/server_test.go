package chat

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func TestServer_MaxClientsLimit(t *testing.T) {
	// 1. Dynamically pick an open port for isolated testing
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	_, port, _ := net.SplitHostPort(listener.Addr().String())
	listener.Close()

	// 2. Run Server asynchronously
	go Start(port)

	// Robust wait loop for TCP server port binding confirmation
	var initialConn net.Conn
	for i := 0; i < 10; i++ {
		initialConn, err = net.Dial("tcp", "127.0.0.1:"+port)
		if err == nil {
			initialConn.Close()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("server failed to start on port %s: %v", port, err)
	}

	// 3. Establish connections up to the maximum limit
	conns := make([]net.Conn, maxClients)
	scanners := make([]*bufio.Scanner, maxClients)
	for i := 0; i < maxClients; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:"+port)
		if err != nil {
			t.Fatalf("failed to connect client %d: %v", i+1, err)
		}
		conns[i] = conn
		t.Cleanup(func() { conn.Close() })

		// Send name BEFORE reading anything
		_, err = conn.Write([]byte(string(rune(65+i)) + "\n"))
		if err != nil {
			t.Fatalf("failed handshake payload for client %d: %v", i+1, err)
		}

		scanner := bufio.NewScanner(conn)
		scanner.Buffer(make([]byte, 0, 2*1024), 256*1024)
		scanners[i] = scanner

		joinMsg := ""
		foundJoin := false
		for scanner.Scan() {
			line := scanner.Text()
			joinMsg += line + "\n"
			if strings.Contains(line, "has joined our chat.") {
				foundJoin = true
				break
			}
		}
		if !foundJoin {
			t.Fatalf("client %d did not receive join confirmation", i+1)
		}
	}

	// Ensure all clients are fully registered before testing overflow
	time.Sleep(200 * time.Millisecond)

	// 4. Attempt the overflow connection with an evaluation loop
	var overflowConn net.Conn
	var response string

	// Poll briefly to let the asynchronous map registration complete safely
	for attempts := 0; attempts < 10; attempts++ {
		overflowConn, err = net.Dial("tcp", "127.0.0.1:"+port)
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}

		scanner := bufio.NewScanner(overflowConn)
		if scanner.Scan() {
			response = strings.TrimSpace(scanner.Text())
			overflowConn.Close()
			break
		}
		overflowConn.Close()
		time.Sleep(10 * time.Millisecond)
	}

	// 5. Final assertion validation
	expected := "Chat is full. Try again later."
	if response != expected {
		t.Errorf("expected rejection message %q, got %q", expected, response)
	}
}