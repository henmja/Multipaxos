// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*DistributedExtCl]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *DistributedExtCl

	// Unregister requests from clients.
	unregister chan *DistributedExtCl
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *DistributedExtCl),
		unregister: make(chan *DistributedExtCl),
		clients:    make(map[*DistributedExtCl]bool),
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
