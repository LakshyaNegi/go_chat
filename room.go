package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	name    string
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

func newRoom(name string) *room {
	return &room{
		name:    name,
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// vars := mux.Vars(req)
	// roomName, _ := vars["roomName"]
	// r.name = roomName

	fmt.Println(r.name)

	socket, err := upgrader.Upgrade(rw, req, nil)
	//socket := *websocket.Conn
	if err != nil {
		log.Fatal("ServeHTTP : ", err.Error)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	client := &client{socket: socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		username: authCookie.Value,
	}

	r.join <- client
	defer func() { r.leave <- client }()

	go client.read()
	client.write()
}
