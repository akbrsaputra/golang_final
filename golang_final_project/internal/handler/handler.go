package handler

import (
	"sync"

	"github.com/akbrsaputra/golang_final_project/chatapp/internal/client"
	"github.com/akbrsaputra/golang_final_project/chatapp/internal/message"
)

// MessageHandler merupakan antarmuka untuk menangani pesan
type MessageHandler interface {
	HandleMessage(client *client.Client, msg message.Message)
}

// MessageProcessor merupakan implementasi default dari antarmuka MessageHandler
type MessageProcessor struct {
	clientsMutex sync.Mutex
	clients      map[*client.Client]bool
	broadcast    chan message.Message
}

// NewMessageProcessor membuat instance MessageProcessor baru
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		clients:   make(map[*client.Client]bool),
		broadcast: make(chan message.Message),
	}
}

// StartBroadcasting mulai menyiarkan pesan ke semua klien
func (m *MessageProcessor) StartBroadcasting() {
	defer close(m.broadcast)

	for msg := range m.broadcast {
		m.clientsMutex.Lock()
		for client := range m.clients {
			go func(client *client.Client, msg message.Message) {
				client.SendMessage(msg)
			}(client, msg)
		}
		m.clientsMutex.Unlock()
	}
}

// HandleMessage mengimplementasikan antarmuka MessageHandler
func (m *MessageProcessor) HandleMessage(client *client.Client, msg message.Message) {
	m.broadcast <- msg
}

// RegisterClient mendaftarkan klien ke dalam MessageProcessor
func (m *MessageProcessor) RegisterClient(client *client.Client) {
	m.clientsMutex.Lock()
	m.clients[client] = true
	m.clientsMutex.Unlock()
}

// UnregisterClient melepaskan klien dari MessageProcessor
func (m *MessageProcessor) UnregisterClient(client *client.Client) {
	m.clientsMutex.Lock()
	delete(m.clients, client)
	m.clientsMutex.Unlock()
}
