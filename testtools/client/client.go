package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("Connected to the Telnet server...")

	// Send WILL and DO commands
	sendTelnetCommand(conn, 251, 1)  // WILL ECHO
	sendTelnetCommand(conn, 253, 31) // DO NAWS

	// Sleep for a moment to wait for server responses
	time.Sleep(2 * time.Second)

	// Send WONT and DONT commands
	sendTelnetCommand(conn, 252, 1)  // WONT ECHO
	sendTelnetCommand(conn, 254, 31) // DONT NAWS

	// Send plain text message to the server
	message := "Hello, Telnet Server!\n"
	conn.Write([]byte(message))
	fmt.Printf("Sent message: %s\n", message)

	// Read and print server response
	reader := bufio.NewReader(conn)
	response, _ := reader.ReadString('\n')
	fmt.Printf("Received message: %s", response)
}

func sendTelnetCommand(conn net.Conn, command byte, option byte) {
	data := []byte{255, command, option}
	conn.Write(data)
	fmt.Printf("Sent: %v\n", data)
}
