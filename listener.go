// Package: telnetter
// Description: Provides a set of tools for building and managing a telnet server. 
// Supports various Telnet commands, options, and subnegotiations.
// Git Repository: [URL not provided in the source]
// License: [License type not provided in the source]
package telnetter

import (
	"net"
	"time"
)

// Telnet commands, states, and options constants
const (
	// READ_STATE_NORMAL is the default state when reading data from the connection.
	READ_STATE_NORMAL  = 1
	// READ_STATE_COMMAND is the state when reading a Telnet command.
	READ_STATE_COMMAND = 2
	// READ_STATE_SUBNEG is the state when reading a Telnet subnegotiation.
	READ_STATE_SUBNEG  = 3

	// Telnet commands
	// TN_INTERPRET_AS_COMMAND is the command to interpret the next byte as a Telnet command.
	TN_INTERPRET_AS_COMMAND = 255
	// TN_ARE_YOU_THERE is the command to request a response from the server.
	TN_ARE_YOU_THERE        = 246

	// Telnet command responses
	//TN_WILL is the response to a request to enable an option.
	TN_WILL                 = 251
	// TN_WONT is the response to a request to disable an option.
	TN_WONT                 = 252
	// TN_DO is the response to a request to enable an option.
	TN_DO                   = 253
	// TN_DONT is the response to a request to disable an option.
	TN_DONT                 = 254
	// TN_SUBNEGOTIATION_START is the response to a request to start a subnegotiation.
	TN_SUBNEGOTIATION_START = 250
	// TN_SUBNEGOTIATION_END is the response to a request to end a subnegotiation.
	TN_SUBNEGOTIATION_END   = 240

	// Telnet options
	//TN_ECHO is the option to enable echo.
	TN_ECHO              = 1
	// TN_SUPPRESS_GO_AHEAD is the option to suppress Go Ahead.
	TN_SUPPRESS_GO_AHEAD = 3
	// TN_STATUS is the option to request the status of all active options.
	TN_STATUS            = 5
	// TN_TERMINAL_TYPE is the option to request the terminal type.
	TN_TERMINAL_TYPE     = 24
	// TN_WINDOW_SIZE is the option to request the window size.
	TN_WINDOW_SIZE       = 31
	// TN_BINARY is the option to request binary data.
	TN_BINARY            = 0
	// TN_TIMING_MARK is the option to request a timing mark.
	TN_TIMING_MARK       = 6
)

// Listener represents a Telnet server listener.
type Listener struct {
	listener net.Listener
	timeout  time.Duration
}


// Title: Start Telnet Listener
// Description: Initializes and starts a new Telnet server listener on the provided bind address.
// Function: func Listen(bind string) (*Listener, error)
// CalledWith: listener, err := Listen("localhost:23")
// ExpectedOutput: A Listener instance or an error if there's an issue starting the listener.
// Example: telnetListener, err := Listen("0.0.0.0:23")
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


// Title: Set Connection Timeout
// Description: Defines the maximum duration the server will wait for data on a connection.
// Function: func (l *Listener) SetTimeout(dur time.Duration)
// CalledWith: listener.SetTimeout(5 * time.Second)
// ExpectedOutput: None, sets the timeout duration for connections.
// Example: listener.SetTimeout(10 * time.Second)
func (l *Listener) SetTimeout(dur time.Duration) {
	l.timeout = dur
}


// Title: Accept Connection
// Description: Waits for and returns the next connection to the listener.
// Function: func (l *Listener) Accept() (*Conn, error)
// CalledWith: conn, err := listener.Accept()
// ExpectedOutput: The next connection to the listener or an error if there's an issue.
// Example: connection, error := listener.Accept()
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

// handleConnection manages the entire lifecycle of a telnet connection. 
// It reads data, processes Telnet commands, and manages subnegotiations.
func handleConnection(conn net.Conn, timeout time.Duration, telCon *Conn) {
    // Create a session for the new connection using the session manager from session.go
    sm := NewSessionManager()
    _ = sm.CreateSession(telCon)

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
					handleSubnegotiation(subnegData, conn, telCon)
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

// handleWillOption processes a received WILL option from the client.
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
	case TN_BINARY:
		// If the client wants to send binary data, acknowledge it
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DO, TN_BINARY})
	default:
		// For all other options, decline the client's request
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DONT, option})
	}
}

// handleWontOption processes a received WONT option from the client.
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

// handleDoOption processes a received DO option from the client.
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
	case TN_BINARY:
		// If the client wants the server to send binary data, acknowledge it
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WILL, TN_BINARY})
	case TN_TIMING_MARK:
		// Respond immediately to a timing mark request to synchronize with the client
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_TIMING_MARK})
	default:
		// For other options, we'll just decline for now
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, option})
	}
}

// handleSubnegotiation deals with received subnegotiation data from the client.
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

// handleDontOption processes a received DONT option from the client.
func handleDontOption(option byte, conn net.Conn) {
	switch option {
	case TN_ECHO:
		// If the client doesn't want the server to echo, disable echo on the server side
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, TN_ECHO})
	case TN_SUPPRESS_GO_AHEAD:
		// Suppress Go Ahead on the server side
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_DO, TN_SUPPRESS_GO_AHEAD})
	case TN_STATUS, TN_TERMINAL_TYPE, TN_WINDOW_SIZE:
		// For these options, we'll acknowledge the client's request to not perform the option
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, option})
	default:
		// For all other options, just acknowledge the client's request
		conn.Write([]byte{TN_INTERPRET_AS_COMMAND, TN_WONT, option})
	}
}