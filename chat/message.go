package chat

import (
	"encoding/json"
	"fmt"
	"time"
)

// different types of messages sent over the websocket
// we can execute different actions according to the type
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const ( // types of events

	// when new message is sent to the chat
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventChatChange  = "change_room"
)

// when payload is a message sent in chat from client
type SendMessageEvent struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

// message to be broadcasted to clients
type NewMessage struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type changeChatEvent struct {
	Name string `json:"name"`
}

func SendMessageHandler(e Event, c *Client) error {

	var incomingEvent SendMessageEvent
	// extracting contents of the message
	err := json.Unmarshal(e.Payload, &incomingEvent)
	if err != nil {
		return fmt.Errorf("invalid payload in request: %v", err)
	}

	// message to be sent to all clients
	var broadcastMessage NewMessage
	broadcastMessage.Message = incomingEvent.Message
	broadcastMessage.Username = incomingEvent.Username
	broadcastMessage.Sent = time.Now()

	data, err := json.Marshal(broadcastMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Type = EventNewMessage
	outgoingEvent.Payload = data
	fmt.Println("message broadcasted")

	for client := range c.Hub.clients {
		client.broadcast <- outgoingEvent
	}
	return nil
}
func ChatChangeHandler(e Event, c *Client) error {
	var newchat changeChatEvent
	err := json.Unmarshal(e.Payload, &newchat)
	if err != nil {
		return fmt.Errorf("invalid payload in requests: %v ", err)
	}
	c.chatroom = newchat.Name
	return nil
}
