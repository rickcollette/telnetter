package telnetter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type MessageHandler func(*Conn, string)
type DisconnectHandler func(*Conn)
type OptionCallback func(*Conn, byte, bool)

type Conn struct {
    ReadWrite       io.ReadWriter
    Connection      net.Conn
    TerminalType    string
    optionState     map[byte]bool
    optionCallbacks map[byte]OptionCallback
    Lock            sync.Mutex
    MessageHandler  MessageHandler
    DisconnectHandler DisconnectHandler
}

func NewConn(conn net.Conn) *Conn {
    return &Conn{
        ReadWrite:       conn,
        Connection:      conn,
        optionState:     make(map[byte]bool),
        optionCallbacks: make(map[byte]OptionCallback),
    }
}

type DataReader struct {
	sourced  io.Reader
	buffered *bufio.Reader
}

func bufferedDataReader(r io.Reader) *DataReader {
	buffered := bufio.NewReader(r)
	return &DataReader{sourced: r, buffered: buffered}
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

func (c *Conn) Write(b []byte) (int, error) {
    escapedData := make([]byte, 0, len(b)*2)
    for _, byteValue := range b {
        if byteValue == 0xFF {
            escapedData = append(escapedData, 0xFF, 0xFF)
        } else {
            escapedData = append(escapedData, byteValue)
        }
    }
    return c.Connection.Write(escapedData)
}

func (r *DataReader) Read(data []byte) (n int, err error) {
    return r.buffered.Read(data)
}

func (c *Conn) SetOptionCallback(option byte, callback OptionCallback) {
    c.optionCallbacks[option] = callback
}

func (c *Conn) Read(b []byte) (int, error) {
    reader := bufferedDataReader(c.Connection)
    n, err := reader.Read(b)
    if err != nil {
        return n, err
    }
    
    outIdx := 0
    i := 0
    for i < n {
        if b[i] == 0xFF {
            if i+1 < n {
                switch b[i+1] {
                case 250:  // SB
                    endIdx := i + 2
                    for endIdx+1 < n && !(b[endIdx] == 0xFF && b[endIdx+1] == 240) {
                        endIdx++
                    }
                    fmt.Printf("Sub-negotiation Data: %v\\n", b[i+2:endIdx])
                    i = endIdx + 1
                    continue
                case 251:  // WILL
                    c.optionState[b[i+2]] = true
                    if callback, exists := c.optionCallbacks[b[i+2]]; exists {
                        callback(c, b[i+2], true)
                    }
                    c.Write([]byte{0xFF, 253, b[i+2]})
                case 252:  // WONT
                    c.optionState[b[i+2]] = false
                    if callback, exists := c.optionCallbacks[b[i+2]]; exists {
                        callback(c, b[i+2], false)
                    }
                    c.Write([]byte{0xFF, 254, b[i+2]})
                case 253:  // DO
                    c.optionState[b[i+2]] = true
                    if callback, exists := c.optionCallbacks[b[i+2]]; exists {
                        callback(c, b[i+2], true)
                    }
                    c.Write([]byte{0xFF, 251, b[i+2]})
                case 254:  // DONT
                    c.optionState[b[i+2]] = false
                    if callback, exists := c.optionCallbacks[b[i+2]]; exists {
                        callback(c, b[i+2], false)
                    }
                    c.Write([]byte{0xFF, 252, b[i+2]})
                case 0xFF:  // Escaped IAC in data
                    b[outIdx] = b[i]
                    outIdx++
                    i += 2
                    continue
                }
                i += 2
            }
        } else {
            b[outIdx] = b[i]
            outIdx++
            i++
        }
    }
    return outIdx, nil
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
	b := make([]byte, 1)
	_, err := c.Connection.Read(b)
	if err != nil {
		return 0, fmt.Errorf("ReadByte: Error reading byte: %w", err)
	}

	if b[0] == 0xFF {
		command := make([]byte, 2)
		_, err := c.Connection.Read(command)
		if err != nil {
			return 0, fmt.Errorf("error reading command sequence: %w", err)
		}
		return c.ReadByte()
	}

	return b[0], nil
}

func (c *Conn) LocalAddr() net.Addr {
    return c.Connection.LocalAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
    return c.Connection.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
    return c.Connection.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
    return c.Connection.SetWriteDeadline(t)
}