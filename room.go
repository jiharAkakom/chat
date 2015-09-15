package main

type room struct {
	//forward is a channel that hold incoming messages
	//that shold be forwarded to the other clients
	forward chan []byte
}
