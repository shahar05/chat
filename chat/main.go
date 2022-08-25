package main

import (
	"log"
	"net"
)

func main() {
	log.Println("starting...")
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server due to: %s", err.Error())
	}

	defer listener.Close()

	log.Println("started server")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection due to: %s", err.Error())
			continue
		}

		go s.handleConnection(conn)
	}

}
