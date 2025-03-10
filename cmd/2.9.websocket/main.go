package main

import (
	"backend-engineering/pkg/httpx/middleware"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = new(sync.Map)

type WSMessage struct {
	msg []byte
	mt  int
}

type handler struct {
	upgrader  websocket.Upgrader
	broadcast chan WSMessage
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error while upgrading", "error", err.Error())
	}

	defer conn.Close()

	slog.Info("user connected", "conn", conn.RemoteAddr().String())

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("read:", "error", err)
			break
		}

		slog.Info("recv", "msg", message)

		if err := conn.WriteMessage(mt, message); err != nil {
			slog.Error("write:", "error", err)
			break
		}
	}
}

func (h handler) Chat(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("error while upgrading", "error", err.Error())
	}

	defer conn.Close()

	addr := conn.RemoteAddr().String()

	clients.Store(addr, conn)
	slog.Info("user connected", "conn", addr)

	for {
		mt, message, err := conn.ReadMessage()

		if err != nil {
			slog.Error("read:", "error", err)
			clients.Delete(addr)
			return
		}
		slog.Info("recv", "msg", message)

		h.broadcast <- WSMessage{[]byte(fmt.Sprintf("user %v => %s \n", addr, message)), mt}
	}

}

func (h handler) distributeMessages() {
	for {
		msg := <-h.broadcast

		clients.Range(func(k, v any) bool {
			conn, ok := v.(*websocket.Conn)
			if !ok {
				slog.Error("cannot type assert client")

				clients.Delete(k)
				return true
			}

			if err := conn.WriteMessage(msg.mt, msg.msg); err != nil {
				slog.Error("error on writing to socket", "error", err.Error())

				clients.Delete(k)
				return true
			}

			return true
		})
	}
}

func main() {
	var broadcast = make(chan WSMessage)

	h := handler{
		upgrader:  websocket.Upgrader{},
		broadcast: broadcast,
	}

	http.Handle("/echo", middleware.Handler(h))
	http.HandleFunc("/chat", middleware.Func(h.Chat))

	go h.distributeMessages()

	slog.Info("server started", "port", "5432")
	slog.Error(http.ListenAndServe(":5432", nil).Error())
}
