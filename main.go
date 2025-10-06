package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection or client failed to connect")
		}
		go handleConnectons(conn)
	}

}

func handleConnectons(conn net.Conn) {
	defer conn.Close()

	b := make([]byte, 1024)
	reader := bufio.NewReader(conn)
	reader.Read(b)

	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
