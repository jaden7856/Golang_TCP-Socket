package main

import (
	"io"
	"log"
	"net"

	"github.com/hashicorp/yamux"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("TCP accept: %s", err)
		}
		go handle(conn)
	}

}

func handle(conn net.Conn) {

	log.Printf("TCP accepted")

	session, err := yamux.Server(conn, nil)
	if err != nil {
		log.Fatalf("Yamux server: %s", err)
	}

	for {
		sconn, err := session.Accept()
		if err != nil {
			if session.IsClosed() {
				log.Printf("TCP closed")
				break
			}
			log.Printf("Yamux accept: %s", err)
			continue
		}
		go streamServer(sconn)
	}
}

func streamServer(sconn net.Conn) {
	buff := make([]byte, 0xff)
	for {
		n, err := sconn.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Stream read error: %s", err)
			break
		}
		log.Printf("stream sent %d bytes: %s", n, buff[:n])
	}
}
