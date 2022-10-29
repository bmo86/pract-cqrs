package main

import "github.com/gorilla/websocket"

type Client struct {
	Hub     *Hub
	Id      string
	Socket  *websocket.Conn
	OutBand chan []byte
}

func NewClient(h *Hub, socket *websocket.Conn) *Client {
	return &Client{
		Hub:     h,
		Socket:  socket,
		OutBand: make(chan []byte),
	}
}

func (c *Client) Write() {
	for {
		select {
		case msg, ok := <-c.OutBand:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, msg)

		}
	}
}

func (c *Client) Close() {
	c.Socket.Close()
	close(c.OutBand)
}
