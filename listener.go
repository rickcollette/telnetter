package telnetter

import (
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	ReadStateNormal  = 1
	ReadStateCommand = 2
	ReadStateSubnegotiation  = 3

	TNInterpretAsCommand = 255
	TNAreYouThere        = 246

	TNWill                 = 251
	TNWont                 = 252
	TNDo                   = 253
	TNDont                 = 254
	TNSubnegotiationStart = 250
	TNSubnegotiationEnd   = 240

	TNEcho              = 1
	TNSuppressGoAhead = 3
	TNStatus            = 5
	TNTerminalType     = 24
	TNWindowSize       = 31
	TNBinary            = 0
	TNTimingMark       = 6
)

type Listener struct {
	listener net.Listener
	shutdownCh chan struct{}
	timeout  time.Duration
}

func Listen(bind string) (*Listener, error) {
	l, err := net.Listen("tcp4", bind)
	if err != nil {
		return nil, err
	}

	return &Listener{
		listener: l,
		shutdownCh: make(chan struct{}),
		timeout:  0,
	}, nil
}



func (l *Listener) SetTimeout(dur time.Duration) {
	l.timeout = dur
}



func (l *Listener) Accept() (*Conn, error) {
    select {
    case <-l.shutdownCh:
        return nil, errors.New("Listener is shutting down")
    default:
    }

	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	telnetConn := &Conn{
		Connection: conn,
		ReadWrite: conn,
	}

	go handleConnection(conn, l.timeout, telnetConn)

	return telnetConn, nil
}


func handleIACSequence(conn net.Conn) (byte, error) {
    command := make([]byte, 2)
    _, err := conn.Read(command)
    if err != nil {
        return 0, fmt.Errorf("handleIACSequence: Error reading command sequence: %w", err)
    }
    
    switch command[0] {
    case TNEcho:
    case TNSuppressGoAhead:
    default:
        return 0, fmt.Errorf("unhandled Telnet command: %d", command[0])
    }
    
    return 0, nil
}

func handleConnection(conn net.Conn, timeout time.Duration, telnetConn *Conn) {
    sm := NewSessionManager()
    _ = sm.CreateSession(telnetConn)

    state := ReadStateNormal
    echoEnabled := false
    buf := make([]byte, 2048)

    for {
        if timeout > 0 {
            _ = conn.SetReadDeadline(time.Now().Add(timeout))
        }

        read, err := conn.Read(buf)
        if err != nil {
            break
        }

        for i := 0; i < read; i++ {
            switch state {
            case ReadStateNormal:
                if buf[i] == TNInterpretAsCommand {
                    state = ReadStateCommand
                } else if echoEnabled {
                    conn.Write(buf[i : i+1])
                }
            case ReadStateCommand:
                if buf[i] == TNInterpretAsCommand {
                    handleIACSequence(conn)
                }
            }
        }
    }

    _ = conn.Close()

    if telnetConn.DisconnectHandler != nil {
        telnetConn.DisconnectHandler(telnetConn)
    }
}

func (l *Listener) Shutdown() error {
    close(l.shutdownCh)
    return l.listener.Close()
}