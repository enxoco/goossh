package main

import (
	"encoding/base64"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	keats      = "O soft embalmer of the still midnight! Shutting with careful fingers and benign Our gloom-pleased eyes, embower’d from the light, Enshaded in forgetfulness divine; O soothest Sleep! if so it please thee, close, In midst of this thine hymn, my willing eyes, Or wait the amen, ere thy poppy throws Around my bed its lulling charities."
	lineLength = 64
	delay      = 10 * time.Second
	maxConns   = 4
)

func main() {
	ln, err := net.Listen("tcp", ":2222")
	if err != nil {
		log.Panic("Could not create listener:", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":80", nil)

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
		if conns >= maxConns {
			log.Println("Connection attempted but hit max, closing connection")
			conn.Close()
			continue
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
		ln := generateLine(keats)
		_, err := conn.Write(ln)
		if err != nil {
			done <- true
			return
		}
		time.Sleep(delay)
	}
}

func generateLine(text string) []byte {
	length := len(text)
	index := rand.Intn(length - lineLength - 1)
	runes := []rune(keats)
	msg := string(runes[index:index+lineLength]) + "\n\n\n"
	encoded := base64.StdEncoding.EncodeToString([]byte(msg))
	return []byte(encoded)
}

func timeTaken(t time.Time, client string) {
	elapsed := time.Since(t)
	log.Printf("Client disconnected: %s\t\t%s\n", client, elapsed)
}
