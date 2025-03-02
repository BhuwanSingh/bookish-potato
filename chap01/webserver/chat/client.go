package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte // the send channel is to store the messages that are sent to this specific client.
	room   *room
}

func (c *client) read() { // the read the message that the client has sent to the server.
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage() // reading the message on the websocket that is sent by the web client.
		if err != nil {
			return
		}
		c.room.forward <- msg // sending the message to the room on which there are other people.
	}
}

func (c *client) write() { // to read the messages that others have written.
	defer c.socket.Close()
	for msg := range c.send {
		// write the message on the websocket of the front-end client that can be read by the person
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
