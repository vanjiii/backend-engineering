package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	// Start listening on port 9000
	listener, err := net.Listen("tcp", ":3001")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 3001")

	for {
		// Accept a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		// Read client message until newline
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Client %s disconnected\n", conn.RemoteAddr())
			break
		}

		fmt.Printf("Received: %s", message)

		// Echo message back to client
		response := "Echo: " + message
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing back to client:", err)
			break
		}
	}
}
