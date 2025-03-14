package main

import (
	"log"
	"os"
)

func main() {
	cmd := os.Args[1]
	if cmd == "" {
		log.Fatal("missing command; aborting")
	}

	switch cmd {
	case "send":
		send()
	case "recv":
		recv()
	default:
		log.Fatalf("unrecognized command: %v", cmd)
	}

}
