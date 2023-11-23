package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/akbrsaputra/golang_final_project/internal/client"
	"github.com/akbrsaputra/golang_final_project/internal/handler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebSocketServer merupakan server WebSocket
type WebSocketServer struct {
	MessageProcessor *handler.MessageProcessor
}

// NewWebSocketServer membuat instance WebSocketServer baru
func NewWebSocketServer(messageProcessor *handler.MessageProcessor) *WebSocketServer {
	return &WebSocketServer{
		MessageProcessor: messageProcessor,
	}
}

func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("Username tidak boleh kosong")
		return
	}

	client := client.NewClient(username, conn, s.MessageProcessor)

	s.MessageProcessor.RegisterClient(client)

	s.MessageProcessor.HandleMessage(nil, message.Message{Username: "Server", Content: username + " bergabung ke dalam chat"})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go client.ListenForMessages(&wg)

	wg.Wait()

	s.MessageProcessor.UnregisterClient(client)

	s.MessageProcessor.HandleMessage(nil, message.Message{Username: "Server", Content: username + " keluar dari chat"})
}

// StartServer memulai server WebSocket
func (s *WebSocketServer) StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", s.handleWebSocket)
	http.Handle("/", r)

	go s.MessageProcessor.StartBroadcasting()

	fmt.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
