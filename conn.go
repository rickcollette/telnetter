// Package: telnetter
// Description: Provides utilities for managing Telnet connections.
// Git Repository: [URL not provided in the source]
// License: [License type not provided in the source]
package telnetter

import (
	"bytes"
	"io"
	"net"
	"sync"
)

// MessageHandler defines the function signature for handling received messages.
// It provides a mechanism for custom logic to be executed upon receiving a message.
type MessageHandler func(*Conn, string)

// DisconnectHandler defines the function signature for handling disconnect events.
// It provides a mechanism for custom logic to be executed upon a client disconnecting.
type DisconnectHandler func(*Conn)

// Conn represents a Telnet connection with utilities for reading/writing data,
// managing terminal type/size, and setting custom handlers for messages and disconnects.
type Conn struct {
	rw         io.ReadWriter
	conn       net.Conn
	terminalType  string
	terminalWidth int
	terminalHeight int
	msgHandler MessageHandler
	disHandler DisconnectHandler
	lock       sync.Mutex
}


// Title: Get Terminal Type
// Description: Retrieves the terminal type associated with the connection.
// Function: func (c *Conn) GetTerminalType() string
// CalledWith: conn.GetTerminalType()
// ExpectedOutput: A string representing the terminal type.
// Example: termType := conn.GetTerminalType()
func (c *Conn) GetTerminalType() string {
	return c.terminalType
}


// Title: Close Connection
// Description: Closes the Telnet connection.
// Function: func (c *Conn) Close() error
// CalledWith: conn.Close()
// ExpectedOutput: An error if the closing fails, otherwise nil.
// Example: err := conn.Close()
func (c *Conn) Close() error {
	return c.conn.Close()
}


// Title: Get Remote Address
// Description: Returns the remote address of the connection.
// Function: func (c *Conn) RemoteAddr() net.Addr
// CalledWith: addr := conn.RemoteAddr()
// ExpectedOutput: The remote address of the connection.
// Example: remoteAddress := conn.RemoteAddr()
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}


// Title: Set Message Handler
// Description: Sets the message handler for the connection.
// Function: func (c *Conn) SetMessageHandler(handler MessageHandler)
// CalledWith: conn.SetMessageHandler(handlerFunction)
// ExpectedOutput: None, sets the message handler for the connection.
// Example: conn.SetMessageHandler(func(conn *Conn, msg string) { ... })
func (c *Conn) SetMessageHandler(handler MessageHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.msgHandler = handler
}


// Title: Set Disconnect Handler
// Description: Sets the disconnect handler for the connection.
// Function: func (c *Conn) SetDisconnectHandler(handler DisconnectHandler)
// CalledWith: conn.SetDisconnectHandler(handlerFunction)
// ExpectedOutput: None, sets the disconnect handler for the connection.
// Example: conn.SetDisconnectHandler(func(conn *Conn) { ... })
func (c *Conn) SetDisconnectHandler(handler DisconnectHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.disHandler = handler
}

// Write writes a slice of bytes to the connection.
func (c *Conn) Write(p []byte) (int, error) {
	return c.rw.Write(p)
}

// Read reads a slice of bytes from the connection.
func (c *Conn) Read(p []byte) (int, error) {
	return c.rw.Read(p)
}

// WriteString writes a string to the connection.
func (c *Conn) WriteString(s string) error {
	_, err := c.Write([]byte(s))
	return err
}

// ReadString reads a string until a newline character is encountered.
func (c *Conn) ReadString() (string, error) {
	var buffer bytes.Buffer
	for {
		b, err := c.ReadByte()
		if err != nil {
			return "", err
		}
		if b == '\n' {
			break
		}
		buffer.WriteByte(b)
	}
	return buffer.String(), nil
}

// ReadByte reads a single byte from the connection.
func (c *Conn) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := c.conn.Read(buf)
	return buf[0], err
}
