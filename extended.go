// Package: telnetter
// Description: Extended utilities for formatting and sending messages in Telnet sessions.
// Git Repository: [URL not provided in the source]
// License: [License type not provided in the source]
package telnetter

import (
	"net"
	"strings"
)


// Title: Send Formatted Message
// Description: Formats and sends a message to the client based on its terminal dimensions.
// Function: func SendMessage(conn net.Conn, message string, width int)
// CalledWith: SendMessage(conn, "Hello, World!", 80)
// ExpectedOutput: None, sends the formatted message to the client.
// Example: SendMessage(conn, "This is a test message.", 50)
func SendMessage(conn net.Conn, message string, width int) {
    lines := formatMessageForTerminal(message, width)
    for _, line := range lines {
        conn.Write([]byte(line + "\r\n"))
    }
}


// Title: Format Message for Terminal
// Description: Breaks down a message into lines that fit within the terminal width.
// Function: func formatMessageForTerminal(message string, width int) []string
// CalledWith: lines := formatMessageForTerminal("This is a long message.", 10)
// ExpectedOutput: A slice of strings representing lines of the formatted message.
// Example: formattedLines := formatMessageForTerminal("Hello, World!", 5)
func formatMessageForTerminal(message string, width int) []string {
    words := strings.Split(message, " ")
    var lines []string
    var currentLine string

    for _, word := range words {
        if len(currentLine)+len(word)+1 > width {
            lines = append(lines, currentLine)
            currentLine = ""
        }
        if currentLine != "" {
            currentLine += " "
        }
        currentLine += word
    }

    if currentLine != "" {
        lines = append(lines, currentLine)
    }

    return lines
}
