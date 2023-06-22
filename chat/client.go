package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
)

type Client struct {
	Conn      *websocket.Conn
	Hub       *Hub
	broadcast chan Event
	chatroom  string
}

func NewClient(c *websocket.Conn, h *Hub, chatroom string) *Client {
	return &Client{
		Conn:      c,
		Hub:       h,
		broadcast: make(chan Event),
	}
}

func (c *Client) readMessages(ctx context.Context, done chan<- int) {
	defer func() {
		c.Hub.RemoveClient(c)
		done <- 1
	}()
	for {
		_, payload, err := c.Conn.Read(ctx)
		if err != nil {
			log.Printf("error reading message:%s \n", err)
			break
		}
		var request Event
		err = json.Unmarshal(payload, &request)
		if err != nil {
			log.Printf("error unmarshling message:%s \n", err)
			break
		}
		println(string(request.Payload))

		err = c.Hub.routeEvent(request, c)
		if err != nil {
			log.Printf("error handling message:%s \n", err)
			break

		}
	}
}

func (c *Client) writeMessages(ctx context.Context) {
	defer func() {
		c.Hub.RemoveClient(c)
	}()
	for {
		fmt.Println("writing messages")
		message, ok := <-c.broadcast
		fmt.Println("message recieved")
		if !ok {
			err := c.Conn.Close(websocket.StatusNormalClosure, "")
			if err != nil {
				log.Println("connection closed: ", err)
			}
			return
		}
		data, err := json.Marshal(message)
		if err != nil {
			log.Println("error marshaling message: ", err)
		}
		err = writeTimeout(ctx, time.Second*5, c.Conn, data)
		if err != nil {
			log.Println("error writing message to client", err)
		}
	}
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
