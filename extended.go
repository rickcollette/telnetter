package telnetter

import (
	"net"
	"strings"
)

// SendMessage formats and sends a message to the client based on its terminal dimensions.
func SendMessage(conn net.Conn, message string, width int) {
    lines := formatMessageForTerminal(message, width)
    for _, line := range lines {
        conn.Write([]byte(line + "\r\n"))
    }
}

// formatMessageForTerminal breaks down a message into lines that fit within the terminal width.
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
