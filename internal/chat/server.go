package chat

import (
	"fmt"
	"net"
	"sync"
)

const maxClients = 10

type Server struct {
	clients  map[*Client]bool
	history  []string
	mutex    sync.Mutex
	messages chan string
}

func NewServer() *Server {
	return &Server{
		clients:  make(map[*Client]bool),
		history:  make([]string, 0),
		messages: make(chan string, 100),
	}
}

func Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on the port :%s\n", port)

	server := NewServer()
	
	// Start message broadcaster
	go server.broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		server.mutex.Lock()
		if len(server.clients) >= maxClients {
			server.mutex.Unlock()
			conn.Write([]byte("Chat is full. Try again later.\n"))
			conn.Close()
			continue
		}
		server.mutex.Unlock()

		go server.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	client := NewClient(conn, s)

	if err := client.welcomeMessage(); err != nil {
		client.Close()
		return
	}

	s.mutex.Lock()
	for existingClient := range s.clients {
		if existingClient.name == client.name {
			s.mutex.Unlock()
			conn.Write([]byte("Name already taken. Disconnecting.\n"))
			client.Close()
			return
		}
	}

	// 1. Register the client safely inside the pool
	s.clients[client] = true
	s.mutex.Unlock()

	// 2. Fire up the client's internal dynamic I/O loops first
	go client.Write()
	go client.Read()

	// 3. Catch up on previous chat rooms history elements safely 
	s.sendHistory(client)

	// 4. Notify all other active loops about the connection event
	joinMsg := FormatSystemMessage(client.name + " has joined our chat.")
	s.messages <- joinMsg

	s.mutex.Lock()
	s.history = append(s.history, joinMsg)
	s.mutex.Unlock()
}

func (s *Server) broadcaster() {
	for msg := range s.messages {
		s.mutex.Lock()
		for client := range s.clients {
			select {
			case client.ch <- msg:
			default:
			}
		}
		s.mutex.Unlock()
	}
}

func (s *Server) broadcast(message string, sender *Client) {
	if message == "" {
		return
	}

	formattedMsg := FormatMessage(sender.name, message)

	s.mutex.Lock()
	s.history = append(s.history, formattedMsg)
	s.mutex.Unlock()

	s.mutex.Lock()
	for client := range s.clients {
		if client != sender {
			select {
			case client.ch <- formattedMsg:
			default:
			}
		}
	}
	s.mutex.Unlock()
}

func (s *Server) removeClient(client *Client) {
	s.mutex.Lock()
	if _, exists := s.clients[client]; exists {
		delete(s.clients, client)
		s.mutex.Unlock()
		
		leaveMsg := FormatSystemMessage(client.name + " has left our chat.")
		s.messages <- leaveMsg
		
		s.mutex.Lock()
		s.history = append(s.history, leaveMsg)
		s.mutex.Unlock()
		
		client.Close()
	} else {
		s.mutex.Unlock()
	}
}

func (s *Server) sendHistory(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	for _, msg := range s.history {
		client.ch <- msg
	}
}