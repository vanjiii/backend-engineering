package main

import (
	"fmt"
	"net"
)

func main() {
	// Set the UDP address to listen on
	addr := net.UDPAddr{
		Port: 3001,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// Listen on the UDP address
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server listening on", addr.String())

	// array of N bytes
	// how much you can handle, datagram wise
	buffer := make([]byte, 1024)

	for {
		// Read data from connection
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			continue
		}

		// Print received message
		message := string(buffer[:n])
		fmt.Printf("Received '%s' from %s\n", message, clientAddr)

		// Echo back the message
		_, err = conn.WriteToUDP([]byte("Echo: "+message), clientAddr)
		if err != nil {
			fmt.Println("Error writing:", err)
		}
	}
}
