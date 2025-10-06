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

	// Continuously accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		// Handle each connection concurrently in a new goroutine
		go handleConnections(conn)
	}
}

func handleConnections(conn net.Conn) {
	// Close the connection when this function returns
	defer conn.Close()

	// Create a buffered reader to read from the connection easily
	reader := bufio.NewReader(conn)

	// Prepare to read all HTTP headers from the connection
	var headers []string

	// Loop to read header lines terminated by '\n'
	for {
		line, err := reader.ReadString('\n')
		fmt.Println("line", line)
		if err != nil {
			fmt.Println("Error reading from connection:", err.Error())
			break
		}
		// Remove trailing carriage return and newline characters
		line = strings.TrimRight(line, "\r\n")
		// Empty line indicates end of headers
		if line == "" {
			break
		}
		// Append header line to headers slice
		headers = append(headers, line)
	}

	// Print all headers received (for debugging)
	fmt.Println("headers", headers)

	// Write a simple HTTP 200 OK response with headers end
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
