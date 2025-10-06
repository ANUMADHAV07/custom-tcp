package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var directory = flag.String("directory", "", "Absolute path to the files directory")

func main() {
	flag.Parse()
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
	method := strings.Split(headers[0], " ")[0]

	// Extract the HTTP request path from the request line (e.g., "GET /path HTTP/1.1")
	requestLine := headers[0]
	parts := strings.Split(requestLine, " ")

	if len(parts) < 2 {
		// Invalid request line fallback response
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}
	path := parts[1]

	fmt.Println("path", path)

	var filesParam string

	if strings.HasPrefix(strings.ToLower(path), "/files/") {
		filesParam = strings.TrimSpace(path[len("/files/"):])
	}
	fmt.Println("fileparam", filesParam)

	var pathParams string
	if strings.HasPrefix(strings.ToLower(path), "/echo/") {
		pathParams = strings.TrimSpace(path[len("/echo/"):])
	}

	// if !checkFile(filesParam) {
	// 	conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	// 	return
	// }

	reqBody := readData(headers, *reader)
	fmt.Println("req", string(reqBody))

	compressData := gzipCompress(pathParams)
	var gzipRes []byte
	if compressData != nil {
		gzipRes = compressData
	}

	if method == "GET" {
		// Handle routing based on the extracted path
		routeRequest(path, filesParam, pathParams, gzipRes, headers, conn)
	}

	if method == "POST" {
		created, err := createFile(filesParam, string(reqBody))
		if err != nil {
			fmt.Println(err.Error())
		}
		if created {
			conn.Write([]byte("HTTP/1.1 201 OK\r\n\r\n"))
		}
	}

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
func routeRequest(path, fileparam, pathParams string, gzipRes []byte, headers []string, conn net.Conn) {
	// Compute content length of the response body based on path length
	contentLength := len(path)
	gzipLen := len(gzipRes)

	switch path {
	case "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	case "/files/" + fileparam:
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", contentLength, fileparam)
		conn.Write([]byte(response))
	case "/echo/" + pathParams:
		if checkCompression(headers) {
			resp := fmt.Sprintf(
				"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%s",
				gzipLen,
				gzipRes,
			)
			conn.Write([]byte(resp))

		} else {
			resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"
			conn.Write([]byte(resp))
		}
	default:
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	}
}

// func checkFile(fileName string) bool {
// 	filePath := filepath.Join(*directory, fileName)
// 	info, err := os.Stat(filePath)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return false
// 	}
// 	fmt.Println(info.Name())
// 	return true
// }

func readData(headers []string, reader bufio.Reader) []byte {
	contentLength := 0
	for _, line := range headers {
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			value := strings.TrimSpace(line[len("content-length:"):])
			contentLength, _ = strconv.Atoi(value)
			break
		}
	}
	reqBody := make([]byte, contentLength)
	reader.Read(reqBody)
	return reqBody
}

func createFile(fileName, data string) (bool, error) {
	filePath := filepath.Join(*directory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	defer file.Close()

	_, err = file.Write([]byte(data))
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	return true, nil
}

func checkCompression(headers []string) bool {
	var compressionSchema string

	for _, line := range headers {
		if strings.HasPrefix(strings.ToLower(line), "accept-encoding:") {
			value := strings.TrimSpace(line[len("Accept-Encoding:"):])
			fmt.Println("value", value)
			compressionSchema = value
		}
	}
	tokens := strings.Split(compressionSchema, ",")
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "gzip" {
			return true
		}
	}
	return false
}

func gzipCompress(data string) []byte {
	if data != "" {
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write([]byte(data)); err != nil {
			log.Fatal(err)
		}
		if err := gz.Close(); err != nil {
			log.Fatal(err)
		}
		return b.Bytes()
	}
	return nil
}
