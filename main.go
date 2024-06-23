package main

import (
	"log"
	"net"
	"time"
)

const (
	delay      = 10 * time.Second
)

func main() {
	// Create a serve on port 2222
	ln, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Panic("Could not create listener:", err)
	}

	var conns int
	connDone := make(chan bool)
	log.Println("Starting listening on 2222")
	go func() {
		for {
			if <-connDone {
				conns--
				log.Println("Number of connections is now", conns)
			}
		}
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic("could not accept connection: ", err)
		}

		conns++
		log.Printf("Connection from %s - Number of connections is now %d", conn.RemoteAddr().String(), conns)
		go handleConnection(conn, connDone)
	}
}

func handleConnection(conn net.Conn, done chan<- bool) {
	defer conn.Close()
	defer timeTaken(time.Now(), conn.RemoteAddr().String())

	for {
		// This is the banner we want to send to the client
		_, err := conn.Write([]byte("connecting..." + "\r\n"))

		if err != nil {
			done <- true
			return
		}
		// Wait for this amount of time before sending the next banner.
		time.Sleep(delay)
	}
}

func timeTaken(t time.Time, client string) {
	elapsed := time.Since(t)
	log.Printf("Client disconnected: %s\t\t%s\n", client, elapsed)
}
