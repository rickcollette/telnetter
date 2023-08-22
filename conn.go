package telnetter

import (
	"bytes"
	"io"
	"net"
	"sync"
)

// MessageHandler defines the function signature for handling received messages.
type MessageHandler func(*Conn, string)

// DisconnectHandler defines the function signature for handling disconnect events.
type DisconnectHandler func(*Conn)

// Conn represents a Telnet connection.
type Conn struct {
	rw         io.ReadWriter
	conn       net.Conn
	msgHandler MessageHandler
	disHandler DisconnectHandler
	lock       sync.Mutex
}

// Close closes the Telnet connection.
func (c *Conn) Close() error {
	return c.conn.Close()
}

// RemoteAddr returns the remote address of the connection.
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetMessageHandler sets the message handler for the connection.
func (c *Conn) SetMessageHandler(handler MessageHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.msgHandler = handler
}

// SetDisconnectHandler sets the disconnect handler for the connection.
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
