package main

import (
	"encoding/json"
	"fmt"
)

// A Hub is a 'room' that connects clients.
type Hub struct {
	Id string

	clients map[*Client]bool

	broadcast chan []byte

	register chan *Client

	unregister chan *Client
}

func newHub(id string) *Hub {
	return &Hub{
		Id:         id,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			clientId := client.Id
			for client := range h.clients {
				msg := []byte("some one join room (Id: " + clientId + ")")
				client.send <- msg
			}
			h.clients[client] = true

		case client := <-h.unregister:
			clientId := client.Id
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			for client := range h.clients {
				msg := []byte("some one left room (Id: " + clientId + ")")
				client.send <- msg
			}

		case userMessage := <-h.broadcast:
			var data map[string][]byte

			json.Unmarshal(userMessage, &data)

			for client := range h.clients {
				// self recieve
				if client.Id == string(data["id"]) {
					continue
				}

				select {
				case client.send <- data["message"]:
				default:
					fmt.Println("Everyone is gone, closing the client")
					close(client.send)
					delete(h.clients, client)

				}
			}
		}

	}
}
