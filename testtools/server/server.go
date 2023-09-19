package main

import (
	"fmt"
	"net"

	"github.com/rickcollette/telnetter"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Telnet server started on port 8080...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Wrap the connection using the telnetter package
	telnetConn := telnetter.NewConn(conn)
	defer telnetConn.Close()

	// Set a message handler
	telnetConn.SetMessageHandler(func(c *telnetter.Conn, msg string) {
		fmt.Println("Received message:", msg)
		response := fmt.Sprintf("Hi! The server received: %s\n", msg)
		c.WriteString(response)
	})

	// Set a disconnect handler
	telnetConn.SetDisconnectHandler(func(c *telnetter.Conn) {
		fmt.Println("Client disconnected:", c.RemoteAddr())
	})
	// Set option callbacks for specific Telnet options
	// For this example, we are only setting up a callback for the "ECHO" option (byte value 1)
	telnetConn.SetOptionCallback(1, func(c *telnetter.Conn, option byte, enabled bool) {
		if enabled {
			fmt.Println("Client enabled ECHO option")
		} else {
			fmt.Println("Client disabled ECHO option")
		}
	})

}
