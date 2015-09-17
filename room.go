package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/onufert/trace"
	"github.com/stretchr/objx"
)

//room room definition
type room struct {
	//forward is a channel that hold incoming messages
	//that shold be forwarded to the other clients
	forward chan *message
	//join is a channel for clients wishing to join the room
	join chan *client
	//leave is a channel for clients wishing to leave this room
	leave chan *client
	//clients holds all the current client in this room
	clients map[*client]bool
	// logger
	tracer trace.Tracer
	//avatar
	avatar Avatar
}

// Newroom Create a new chat room
func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			//joining
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			//leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			//forward message to all clients
			r.tracer.Trace("Message received: ", string(msg.Message))
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("-- sent to client")
					//send the message
				default:
					//failed to send
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace("-- failed to send, cleaned up client")
				}
			}
		}

	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
