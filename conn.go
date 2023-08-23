// Package: telnetter
// Description: Provides utilities for managing Telnet connections.
// Git Repository: [URL not provided in the source]
// License: [License type not provided in the source]
package telnetter

import (
	"bytes"
	"errors"
	"fmt"
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
    if c.rw == nil {
        return 0, errors.New("rw is not initialized")
    }
    return c.rw.Write(p)
}

// Read reads a slice of bytes from the connection 
// and prints out the exact byte values being read for debugging purposes.
func (c *Conn) Read(p []byte) (int, error) {
    n, err := c.rw.Read(p)
    if err != nil {
        return n, err
    }
    
    // Print out the exact byte values for debugging
    for _, b := range p[:n] {
        fmt.Printf("Byte read: %d\n", b)
    }
    
    return n, nil
}



// WriteString writes a string to the connection.
func (c *Conn) WriteString(s string) error {
	_, err := c.Write([]byte(s))
	return err
}

// ReadString reads a string until a newline character or sequence is encountered.
func (c *Conn) ReadString() (string, error) {
    var buffer bytes.Buffer
    var prevByte byte
    for {
        b, err := c.ReadByte()
        if err != nil {
            return "", err
        }
        if b == '\n' {
            // If the previous byte was '\r', we have a CRLF sequence.
            // Remove the '\r' from the buffer.
            if prevByte == '\r' {
                bufLen := buffer.Len()
                if bufLen > 0 {
                    buffer.Truncate(bufLen - 1)
                }
            }
            break
        }
        buffer.WriteByte(b)
        prevByte = b
    }
    return buffer.String(), nil
}


// ReadByte reads a single byte from the connection, handling basic Telnet command sequences 
// and prints out the exact byte values being read for debugging purposes.
func (c *Conn) ReadByte() (byte, error) {
    const (
        IAC  = 255 // Interpret As Command
        DO   = 253
        DONT = 254
        WILL = 251
        WONT = 252
    )
    for {
        b := make([]byte, 1)
        _, err := c.conn.Read(b)
        if err != nil {
            return 0, err
        }
        
        // Print out the exact byte value for debugging
        fmt.Printf("Byte read: %d\n", b[0])

        if b[0] == IAC {
            // Read the next two bytes (command and option) and discard them
            command := make([]byte, 2)
            _, err := c.conn.Read(command)
            if err != nil {
                return 0, err
            }
            continue // Go back to the start of the loop to read the next byte
        }
        return b[0], nil
    }
}


