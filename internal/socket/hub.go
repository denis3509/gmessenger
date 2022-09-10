// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"encoding/json"

	"fmt"
	"messenger/pkg/log"

	"messenger/internal/config"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type Payload interface {
	ToJSON() string
}
type IncomingEvent string
type OutgoingEvent string

type EventHandler struct {
	HandleFunc func(*Client, []byte)
	NeedAuth   bool
}

type SocketMessage struct {
	Event   string `json:"event"`
	Payload string `json:"payload"`
}

type Hub struct {

	// Registered clients.
	clients map[*Client]bool
	// Auth clients
	authClients map[int][]*Client
	// event handlers
	eventHandlers map[string][]EventHandler
	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	cfg config.Config
	log  log.Logger
}

func (h *Hub) AddHandler(event string, handler EventHandler) {
	handlers, ok := h.eventHandlers[event]
	if ok {
		handlers = append(handlers, handler)
		h.eventHandlers[event] = handlers
	} else {
		handlers := make([]EventHandler, 1)
		handlers = append(handlers, handler)
		h.eventHandlers[event] = handlers
	}
}

type ErrorPayload struct {
	Message string `json:"message"`
}

func ErrMessage(message string) SocketMessage {

	payload := fmt.Sprintf("{\"message\":\"%s\"}", message)

	return SocketMessage{
		Event:   "error",
		Payload: payload,
	}

}

func (h *Hub) HandleMessage(client *Client, message []byte) {
	var sm SocketMessage
	json.Unmarshal(message, &sm)
	if handlers, ok := h.eventHandlers[sm.Event]; ok {
		for _, h := range handlers {
			if h.NeedAuth {
				// TODO Check Auth
				h.HandleFunc(client, []byte(sm.Payload))
			} else {
				h.HandleFunc(client, []byte(sm.Payload))
			}
		}
	} else { 
		errSm := ErrMessage(fmt.Sprintf("unknown event '%s'", sm.Event))
		client.Send(errSm)
	}
}

// func (h *Hub) ClientsByUserId(userId int) *[]Client {
// 	clients, ok := h.authClients[userId]
// 	if ok {
// 		return clients
// 	} else {
// 		return nil
// 	}
// }

func NewHub(log log.Logger) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		log: log,
	}
}

func (h *Hub) Run() {
	h.log.Info("Running hub")
	for {
		select {
		case client := <-h.register:
			h.log.Info("new client", client.conn.RemoteAddr())
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.deleteClient(client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) deleteClient(client *Client) {
	userId := client.userId

	delete(h.clients, client)
	close(client.send)
	client.conn.Close()
	if client.userId > 0 {
		clients, ok := h.authClients[userId]
		if !ok {
			h.log.Fatalf("client %s not found by id", userId)
		}
		newClients := make([]*Client, len(clients)-1)
		for _, c := range clients {
			if c != client {
				newClients = append(newClients, c)
			}
		}
		h.authClients[userId] = newClients
	}
}

func (h *Hub) SendByClient(c *Client, msg SocketMessage) error {
	if _, ok := h.clients[c]; ok {
		if err := c.Send(msg); err != nil {
			return err
		}

		if c.userId > 0 {
			if err := h.SendByUserId(c.userId, msg); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Hub) SendByUserId(userId int, msg SocketMessage) error {
	if clients, ok := h.authClients[userId]; ok {
		for _, c := range clients {
			if err := c.Send(msg); err != nil {
				h.log.Fatal(err)
				// return err
			}
		}
	}
	return nil
}

func (h *Hub) AddAuthClient(userId int, client *Client) {
	client.userId = userId
	// if clients, ok := h.authClients[userId]; ok {
	// 	clients = append(clients, client)
	// 	h.authClients[userId] = clients
	// } else {
	// 	clients := make([]*Client, 0)
	// 	clients = append(clients, client)
	// 	h.authClients[userId] = clients
	// }
	clients, ok := h.authClients[userId]
	if !ok {
		clients = make([]*Client, 0)
	}
	clients = append(clients, client)
	h.authClients[userId] = clients
}

func (h *Hub) CloseAll() {
	fmt.Println("closing all connections")
	for c := range h.clients {

		h.deleteClient(c)  
	}
}