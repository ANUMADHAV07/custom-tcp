package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// Start listening for TCP connections on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	// Ensure listener is closed when the program exits
	defer listener.Close()

	// Continuously accept incoming TCP connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		// Handle each connection concurrently in its own goroutine
		go handleConnections(conn)
	}
}

// handleConnections processes each accepted TCP connection
func handleConnections(conn net.Conn) {
	// Ensure the connection is closed when the function exits
	defer conn.Close()

	// Wrap connection in a buffered reader to simplify reading by lines
	reader := bufio.NewReader(conn)

	// Read all HTTP request headers from the connection
	headers := readHeaders(reader)

	// Debug: print all headers received
	fmt.Println("headers", headers)

	// Extract the HTTP request path from the request line (e.g., "GET /path HTTP/1.1")
	requestLine := headers[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) < 2 {
		// Invalid request line fallback response
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}
	path := parts[1]

	// Remove leading "/" from path for easier routing
	trimmedPath := strings.TrimPrefix(path, "/")

	// Debug: print extracted paths
	fmt.Println("raw path:", path)
	fmt.Println("trimmed path:", trimmedPath)

	// Handle routing based on the extracted path
	routeRequest(trimmedPath, conn)
}

// readHeaders reads HTTP headers line by line until it encounters an empty line indicating headers completion
func readHeaders(reader *bufio.Reader) []string {
	var headers []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			break
		}
		line = strings.TrimRight(line, "\r\n") // Trim CRLF
		if line == "" {                        // Empty line denotes end of headers
			break
		}
		headers = append(headers, line)
	}
	return headers
}

// routeRequest sends response based on the requested path
func routeRequest(path string, conn net.Conn) {
	// Compute content length of the response body based on path length
	contentLength := len(path)

	// Route based on path values
	if path == "" {
		// Root path "/"
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		// For any other path, respond with path content as plain text
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, path)
		conn.Write([]byte(response))
	}
}
