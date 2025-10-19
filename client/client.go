package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// Client handles TCP communication with the remote server
type Client struct {
	serverAddr string
	timeout    time.Duration
}

// New creates a new Client with the specified server address and timeout
func New(serverAddr string, timeout time.Duration) *Client {
	return &Client{
		serverAddr: serverAddr,
		timeout:    timeout,
	}
}

// SendCommand sends a command to the server and returns the response
func (c *Client) SendCommand(command string) ([]byte, error) {
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	return c.SendMessage(conn, []byte(command))
}

// SendMessage sends a message over an existing connection and returns the response
func (c *Client) SendMessage(conn net.Conn, message []byte) ([]byte, error) {
	conn.SetDeadline(time.Now().Add(c.timeout))

	length := uint16(len(message))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return nil, fmt.Errorf("error writing length: %w", err)
	}

	if _, err := conn.Write(message); err != nil {
		return nil, fmt.Errorf("error writing message: %w", err)
	}

	var responseLength uint16
	if err := binary.Read(conn, binary.BigEndian, &responseLength); err != nil {
		return nil, fmt.Errorf("error reading response length: %w", err)
	}

	response := make([]byte, responseLength)
	if _, err := io.ReadFull(conn, response); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if !bytes.Contains(response, []byte("ok")) {
		return nil, fmt.Errorf("error from remote server: %s", string(response))
	}

	return response, nil
}
