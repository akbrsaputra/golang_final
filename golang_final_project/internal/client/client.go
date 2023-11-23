package client

import (
	"log"
	"sync"

	"github.com/akbrsaputra/golang_final_project/internal/handler"
	"github.com/akbrsaputra/golang_final_project/internal/message"
	"github.com/gorilla/websocket"
)

// Client merupakan struktur untuk merepresentasikan pengguna yang terhubung
type Client struct {
	Username string
	Conn     *websocket.Conn
	Handler  handler.MessageHandler
}

// NewClient membuat instance Client baru
func NewClient(username string, conn *websocket.Conn, handler handler.MessageHandler) *Client {
	return &Client{
		Username: username,
		Conn:     conn,
		Handler:  handler,
	}
}

// SendMessage mengirim pesan ke klien
func (c *Client) SendMessage(msg message.Message) {
	err := c.Conn.WriteJSON(msg)
	if err != nil {
		log.Println(err)
		c.Close()
		c.Handler.(*handler.MessageProcessor).UnregisterClient(c)
	}
}

// Close menutup koneksi klien
func (c *Client) Close() {
	c.Conn.Close()
}

// ListenForMessages mendengarkan pesan dari klien dan menangani dengan handler
func (c *Client) ListenForMessages(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var msg message.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}

		c.Handler.HandleMessage(c, msg)
	}
}
