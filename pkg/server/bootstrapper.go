package server

import (
	"log"
	"net"
	"net/rpc"
)

var (
	listener net.Listener
)

func Run(bindAddr string) {
	// init global variables
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on address:", listener.Addr())

	// Accept connections forever
	rpc.Accept(listener)
}
