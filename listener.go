package telnetter

import (
	"net"
	"time"
)

// Telnet commands, states, and options constants
const (
	READ_STATE_NORMAL  = 1
	READ_STATE_COMMAND = 2
	READ_STATE_SUBNEG  = 3

	TN_INTERPRET_AS_COMMAND = 255
	TN_ARE_YOU_THERE        = 246
	TN_WILL                 = 251
	TN_WONT                 = 252
	TN_DO                   = 253
	TN_DONT                 = 254
	TN_SUBNEGOTIATION_START = 250
	TN_SUBNEGOTIATION_END   = 240

	// Telnet options
	TN_ECHO              = 1
	TN_SUPPRESS_GO_AHEAD = 3
	TN_STATUS            = 5
	TN_TERMINAL_TYPE     = 24
	TN_WINDOW_SIZE       = 31
)

// Listener represents a Telnet server listener.
type Listener struct {
	listener net.Listener
	timeout  time.Duration
}

// Listen starts a new Telnet listener that can accept connections.
func Listen(bind string) (*Listener, error) {
	l, err := net.Listen("tcp4", bind)
	if err != nil {
		return nil, err
	}

	return &Listener{
		listener: l,
		timeout:  0,
	}, nil
}

// SetTimeout sets the duration that the server will wait for data before closing the connection.
func (l *Listener) SetTimeout(dur time.Duration) {
	l.timeout = dur
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (*Conn, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	telCon := &Conn{
		conn: conn,
	}

	go handleConnection(conn, l.timeout, telCon)

	return telCon, nil
}

// handleConnection is a helper function to handle incoming connections and their messages.
func handleConnection(conn net.Conn, timeout time.Duration, telCon *Conn) {
	state := READ_STATE_NORMAL
	echoEnabled := false
	buf := make([]byte, 2048)
	subnegData := make([]byte, 0) // To store bytes related to subnegotiation

	for {
		if int64(timeout) > 0 {
			_ = conn.SetReadDeadline(time.Now().Add(timeout))
		}

		read, err := conn.Read(buf)
		if err != nil {
			break
		}

		for i := 0; i < read; i++ {
			switch state {
			case READ_STATE_NORMAL:
				if buf[i] == TN_INTERPRET_AS_COMMAND {
					state = READ_STATE_COMMAND
				} else if echoEnabled {
					// If echo is enabled, send the character back to the client
					conn.Write(buf[i : i+1])
				}
			case READ_STATE_COMMAND:
				switch buf[i] {
				case TN_WILL:
					i++
					if buf[i] == TN_ECHO {
						echoEnabled = true
					}
					handleWillOption(buf[i], conn)
				case TN_WONT:
					i++
					if buf[i] == TN_ECHO {
						echoEnabled = false
					}
					handleWontOption(buf[i], conn)
				case TN_DO:
					i++
					handleDoOption(buf[i], conn)
				case TN_DONT:
					i++
					handleDontOption(buf[i], conn)
				case TN_SUBNEGOTIATION_START:
					state = READ_STATE_SUBNEG
					subnegData = subnegData[:0] // Reset the slice for new subnegotiation
				default:
					state = READ_STATE_NORMAL
				}
			case READ_STATE_SUBNEG:
				if buf[i] == TN_SUBNEGOTIATION_END {
					state = READ_STATE_NORMAL
					handleSubnegotiation(subnegData, conn, telCon)  // Modified this line to pass telCon
				} else {
					// Append the byte to the subnegotiation data
					subnegData = append(subnegData, buf[i])
				}
			}
		}
	}

	_ = conn.Close()

	if telCon.disHandler != nil {
		telCon.disHandler(telCon)
	}
}

func handleWillOption(option byte, conn net.Conn) {
	switch option {
	case TN_ECHO:
		// Enable echo on the server side
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DO, TN_ECHO})
	case TN_SUPPRESS_GO_AHEAD:
		// Suppress Go Ahead on the server side
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DO, TN_SUPPRESS_GO_AHEAD})
	case TN_TERMINAL_TYPE, TN_WINDOW_SIZE:
		// For these options, we'll just acknowledge the client's request for now
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DO, option})
	default:
		// For all other options, decline the client's request
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DONT, option})
	}
}

func handleWontOption(option byte, conn net.Conn) {
	switch option {
	case TN_ECHO:
		// Disable echo on the server side
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DONT, TN_ECHO})
	case TN_SUPPRESS_GO_AHEAD:
		// Re-enable Go Ahead on the server side (though it's largely obsolete)
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DONT, TN_SUPPRESS_GO_AHEAD})
	default:
		// For all other options, just acknowledge the client's request
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DONT, option})
	}
}

func handleDoOption(option byte, conn net.Conn) {
	switch option {
	case TN_STATUS:
		// Send the status of all active options
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_SUBNEGOTIATION_START, TN_STATUS, TN_WILL, TN_ECHO, TN_WILL, TN_SUPPRESS_GO_AHEAD, TN_SUBNEGOTIATION_END})
	case TN_TERMINAL_TYPE:
		// Respond that we're willing to provide terminal type information
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WILL, TN_TERMINAL_TYPE})
		// Then request the actual terminal type from the client
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_SUBNEGOTIATION_START, TN_TERMINAL_TYPE, 1, TN_SUBNEGOTIATION_END})
	default:
		// For other options, we'll just decline for now
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, option})
	}
}

func handleSubnegotiation(data []byte, conn net.Conn, telCon *Conn) {
	switch data[0] {
	case TN_TERMINAL_TYPE:
		// Store the terminal type reported by the client in the Conn instance
		telCon.terminalType = string(data[2:])
	case TN_WINDOW_SIZE:
		// Decode and store the reported window width and height in the Conn instance
		telCon.terminalWidth = int(data[1])<<8 + int(data[2])
		telCon.terminalHeight = int(data[3])<<8 + int(data[4])
	}
}


func handleDontOption(option byte, conn net.Conn) {
	switch option {
	case TN_ECHO:
		// If the client doesn't want the server to echo, disable echo on the server side
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, TN_ECHO})
	case TN_SUPPRESS_GO_AHEAD:
		// If the client doesn't want to suppress go ahead, acknowledge it
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, TN_SUPPRESS_GO_AHEAD})
	case TN_STATUS, TN_TERMINAL_TYPE, TN_WINDOW_SIZE:
		// For these options, we'll acknowledge the client's request to not perform the option
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, option})
	default:
		// For all other options, just acknowledge the client's request
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, option})
	}
}
