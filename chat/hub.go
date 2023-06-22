package chat

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"nhooyr.io/websocket"
)

type Hub struct {
	clients map[*Client]bool

	sync.RWMutex
	eventHandlers map[string]func(e Event, c *Client) error
}

func NewHub() *Hub {
	h := &Hub{
		clients:       make(map[*Client]bool),
		eventHandlers: make(map[string]func(e Event, c *Client) error),
	}
	h.setupEventHandlers()
	return h
}

func (h *Hub) setupEventHandlers() {
	h.eventHandlers[EventSendMessage] = SendMessageHandler
	h.eventHandlers[EventChatChange] = ChatChangeHandler
}

func (h *Hub) AddClient(client *Client) {
	h.Lock()
	h.clients[client] = true
	defer h.Unlock()
}

func (h *Hub) RemoveClient(client *Client) {
	h.Lock()
	delete(h.clients, client)
	defer h.Unlock()

}

func (h *Hub) routeEvent(e Event, c *Client) error {

	handler, ok := h.eventHandlers[e.Type]
	if ok {
		err := handler(e, c)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("error event not supported")

}
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("New Connection")
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("connection failed: %s\n", err)
		return
	}
	client := NewClient(c, h, "")
	h.AddClient(client)
	done := make(chan int)
	go client.writeMessages(r.Context())
	go client.readMessages(r.Context(), done)
	<-done

}
