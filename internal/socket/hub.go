// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"messenger/internal/config"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.

type Payload interface {
	ToJSON() string
}
type IncomingEvent string
type OutgoingEvent string

const (
	messageCreate      IncomingEvent = "message:create"
	messageList        IncomingEvent = "message:list"
	messageContactList IncomingEvent = "message:contact-list"
	userAuth           IncomingEvent = "user:authenticate"
)

type EventHandler func(*Client, []byte)

type SocketMessage struct {
	Event   string
	Payload string
}

type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Auth clients
	authClients map[int]*Client
	// event handlers
	eventHandlers map[string][]EventHandler
	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	cfg config.Config
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

func (h *Hub) HandleMessage(client *Client, message []byte) {
	var sm SocketMessage
	json.Unmarshal(message, &sm)
	if handlers, ok := h.eventHandlers[sm.Event]; ok {
		for _, h := range handlers {
			h(client, []byte(sm.Payload))
		}
	} else {
		errMsg := ErrorPayload{
			fmt.Sprintf("unknown event '%s'", sm.Event),
		}
		payload, err := json.Marshal(&errMsg)
		if err != nil {
			log.Fatal(err)
		}
		client.Send("error", payload)
	}
}

func (h *Hub) ClientById(userId int) *Client {
	client, ok := h.authClients[userId]
	if ok {
		return client
	} else {
		return nil
	}
}

func newHub() *Hub {
	return &Hub{
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
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
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
