package telnetter

import (
	"net"
	"strings"
)


func SendMessage(conn net.Conn, message string, width int) {
    lines := formatMessageForTerminal(message, width)
    for _, line := range lines {
        conn.Write([]byte(line + "\r\n"))
    }
}


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

func WrapText(text string, width int) string {
    var wrappedText strings.Builder
    words := strings.Fields(text)
    if len(words) == 0 {
        return ""
    }

    wrappedText.WriteString(words[0])
    currentLineLength := len(words[0])
    for _, word := range words[1:] {
        if currentLineLength+len(word)+1 > width {
            wrappedText.WriteString("\n")
            currentLineLength = 0
        } else {
            wrappedText.WriteString(" ")
            currentLineLength++
        }
        wrappedText.WriteString(word)
        currentLineLength += len(word)
    }
    return wrappedText.String()
}

func AlignText(text string, width int, alignment string) string {
    gap := width - len(text)
    if gap <= 0 {
        return text
    }

    switch alignment {
    case "left":
        return text + strings.Repeat(" ", gap)
    case "right":
        return strings.Repeat(" ", gap) + text
    case "center":
        leftPadding := gap / 2
        rightPadding := gap - leftPadding
        return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
    default:
        return text
    }
}

func IndentText(text string, spaces int) string {
    indentation := strings.Repeat(" ", spaces)
    lines := strings.Split(text, "\n")
    for i, line := range lines {
        lines[i] = indentation + line
    }
    return strings.Join(lines, "\n")
}
