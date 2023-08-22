package telnetter

import (
	"net"
	"time"
)

// Telnet commands and states constants
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
	buf := make([]byte, 2048)
	curMsg := ""

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
				} else {
					curMsg += string(buf[i])
				}
			case READ_STATE_COMMAND:
				switch buf[i] {
				case TN_WILL, TN_WONT, TN_DO, TN_DONT:
					// Skipping one byte for the option
					i++
				case TN_SUBNEGOTIATION_START:
					state = READ_STATE_SUBNEG
				default:
					state = READ_STATE_NORMAL
				}
			case READ_STATE_SUBNEG:
				if buf[i] == TN_SUBNEGOTIATION_END {
					state = READ_STATE_NORMAL
				}
			}
		}
	}

	_ = conn.Close()

	if telCon.disHandler != nil {
		telCon.disHandler(telCon)
	}
}
