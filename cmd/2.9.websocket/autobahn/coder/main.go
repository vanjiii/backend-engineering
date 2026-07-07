package main

import (
	"io"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

func echoCopy(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		CompressionMode: websocket.CompressionContextTakeover,
	})

	if err != nil {
		log.Printf("accept error: %v", err)
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "bye")

	c.SetReadLimit(1 << 30) // 1gb

	ctx := r.Context()
	b := make([]byte, 32<<10) // 32 kb
	for {
		typ, reader, err := c.Reader(ctx)
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}

		writer, err := c.Writer(ctx, typ)
		if err != nil {
			log.Printf("writer error: %v", err)
			break
		}

		if typ == websocket.MessageText {
			reader = &validator{r: reader}
		}

		if _, err = io.CopyBuffer(writer, reader, b); err != nil {
			log.Printf("copy error: %v", err)
		}

		if err := writer.Close(); err != nil {
			log.Println("Close:", err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/f", echoCopy)

	log.Println("Coder WebSocket echo server on :9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
