
# Telnetter: A Go Library for Telnet Communication

`telnetter` is a Go library designed to simplify Telnet communication. Whether you're building a Telnet server or client, this library provides a set of tools and abstractions to make the process seamless.

## Features:
- Supports Telnet command interpretation.
- Enhanced error handling for better debugging and problem resolution.
- Modular design for easy extension.

---

## Installation

To include `telnetter` in your Go project, use the following command:

```bash
go get -u github.com/rickcollette/telnetter
```

---

## Usage

### Establishing a Telnet Connection:

```go
import "github.com/rickcollette/telnetter"

// Create a new Telnet connection.
conn, err := telnetter.Dial("localhost:23")
if err != nil {
    log.Fatal(err)
}

// Use the connection to send and receive data.
_, err = conn.WriteString("Hello, Telnet!")
if err != nil {
    log.Fatal(err)
}

// Always close the connection when done.
defer conn.Close()
```

### Handling Telnet Commands:

```go
// Handle incoming Telnet commands.
err = conn.HandleIACCommand()
if err != nil {
    log.Fatal(err)
}
```

---

## Error Handling

`telnetter` uses enhanced error handling to provide context about issues encountered. Always check for errors when using the library's functions and methods. For example:

```go
_, err = conn.WriteString("Hello, Telnet!")
if err != nil {
    // Handle the error accordingly.
    log.Fatalf("Failed to send data: %v", err)
}
```

Errors returned by the library are descriptive and can guide you in resolving potential issues.

---

## Future Improvements

- Additional support for extended Telnet commands.
- Continuous optimization for performance.

---

## Feedback and Contributions

We welcome feedback and contributions to the `telnetter` library. If you encounter any issues or have suggestions, please open an issue or submit a pull request.

---

## License

The `telnetter` library is licensed under the [MIT License](LICENSE).
