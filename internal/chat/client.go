package chat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn    net.Conn
	name    string
	server  *Server
	ch      chan string
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

func NewClient(conn net.Conn, server *Server) *Client {
	return &Client{
		conn:    conn,
		server:  server,
		ch:      make(chan string, 100),
		scanner: bufio.NewScanner(conn),
		writer:  bufio.NewWriter(conn),
	}
}

func (c *Client) Read() {
	for c.scanner.Scan() {
		msg := strings.TrimSpace(c.scanner.Text())
		if msg == "" {
			continue
		}
		c.server.broadcast(msg, c)
	}

	// Client disconnected
	c.server.removeClient(c)
}

func (c *Client) Write() {
	for msg := range c.ch {
		c.writer.WriteString(msg)
		c.writer.Flush()
	}
}

func (c *Client) Close() {
	c.conn.Close()
	close(c.ch)
}

func (c *Client) welcomeMessage() error {
	_, err := c.writer.WriteString(GetWelcomeLogo())
	if err != nil {
		return err
	}
	c.writer.Flush()

	if c.scanner.Scan() {
		name := strings.TrimSpace(c.scanner.Text())
		if name == "" {
			c.writer.WriteString("Name cannot be empty. Disconnecting.\n")
			c.writer.Flush()
			return fmt.Errorf("empty name")
		}
		c.name = name
	}
	return nil
}