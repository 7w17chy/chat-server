package main

import (
	"log"
	"net"
	"time"
)

func main() {
	s := NewServer()
	s.Init()
	go s.Run()

	// FIXME hard-coded number!
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	// FIXME hard-coded number!
	log.Printf("server started on :8888")

	// handle incoming client connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %s", err.Error())
				continue
			}

			go s.NewClient(conn)
		}
	}()

	// keep the server running
	for s.running {
		// sleep for 30 seconds
		time.Sleep(30 * time.Second) // FIXME hard-coded number!
	}
}
