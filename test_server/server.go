package main

import (
	"fmt"
	"net"
	"os"
	"github.com/rickcollette/telnetter" // Assuming telnetter package is in the GOPATH
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Telnet server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Wrap the connection using the telnetter package
	telnetConn := &telnetter.Conn{
		ReadWrite:   conn,
		Connection: conn,
	}

	// Send a welcome message to the client
	telnetConn.Write([]byte("Welcome to the Telnet server!\n"))

	// Continuously read bytes from the client and display them
	for {
		b, err := telnetConn.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			return
		}
		fmt.Printf("Received byte: %d\n", b)
	}
}

