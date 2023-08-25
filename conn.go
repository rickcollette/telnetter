package telnetter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type MessageHandler func(*Conn, string)
type DisconnectHandler func(*Conn)

type Conn struct {
	ReadWrite         io.ReadWriter
	Connection        net.Conn
	TerminalType      string
	TerminalWidth     int
	TerminalHeight    int
	MessageHandler    MessageHandler
	DisconnectHandler DisconnectHandler
	Lock              sync.Mutex
}

type DataReader struct {
	sourced  io.Reader
	buffered *bufio.Reader
}

func bufferedDataReader(r io.Reader) *DataReader {
	buffered := bufio.NewReader(r)

	reader := DataReader{
		sourced:  r,
		buffered: buffered,
	}

	return &reader
}

func (c *Conn) GetTerminalType() string {
	return c.TerminalType
}

func (c *Conn) Close() error {
	return c.Connection.Close()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.Connection.RemoteAddr()
}

func (c *Conn) SetMessageHandler(handler MessageHandler) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.MessageHandler = handler
}

func (c *Conn) SetDisconnectHandler(handler DisconnectHandler) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	c.DisconnectHandler = handler
}

func (c *Conn) Write(p []byte) (int, error) {
	if c.ReadWrite == nil {
		return 0, errors.New("ReadWrite is not initialized")
	}
	return c.ReadWrite.Write(p)
}

func (r *DataReader) Read(data []byte) (n int, err error) {
    return r.buffered.Read(data)
}

func (c *Conn) Read(b []byte) (int, error) {
    reader := bufferedDataReader(c.Connection)
    return reader.Read(b)
}
func (c *Conn) WriteString(s string) error {
	_, err := c.Write([]byte(s))
	return err
}

func (c *Conn) ReadString() (string, error) {
	var buffer bytes.Buffer
	var prevByte byte
	for {
		b, err := c.ReadByte()
		if err != nil {
			return "", err
		}
		if b == '\n' {
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

func (c *Conn) ReadByte() (byte, error) {
	const IAC = 255 // Interpret As Command

	b := make([]byte, 1)
	_, err := c.Connection.Read(b)
	if err != nil {
		return 0, fmt.Errorf("ReadByte: Error reading byte: %w", err)
	}

	if b[0] == IAC {
		_, err := handleIACSequence(c.Connection)
		if err != nil {
			return 0, err
		}
		return c.ReadByte()
	}

	return b[0], nil
}

func (c *Conn) HandleIACCommand() error {
	const (
		DO   = 253
		DONT = 254
		WILL = 251
		WONT = 252
	)
	command := make([]byte, 2)
	_, err := c.Connection.Read(command)
	if err != nil {
		return fmt.Errorf("error reading command sequence: %w", err)
	}
	switch command[0] {
	case DO:
		fmt.Printf("Client requests to DO option %d\n", command[1])
	case DONT:
		fmt.Printf("Client requests to DONT do option %d\n", command[1])
	case WILL:
		fmt.Printf("Client offers to WILL do option %d\n", command[1])
	case WONT:
		fmt.Printf("Client offers to WONT do option %d\n", command[1])
	default:
		fmt.Printf("Unhandled Telnet command: %d\n", command[0])
	}
	return nil
}
