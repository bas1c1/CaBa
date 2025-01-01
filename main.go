package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var maindb db = db{"name", 3}
var cache_ cache
var q *queue = newQueue(8192)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	tr := transaction{
		-1,
		parseRequest(string(buffer[:n])),
	}

	_, res := q.Add(tr)
	s := string(<-res)

	responseStr := fmt.Sprintf("%v", s)
	conn.Write([]byte(responseStr))

	conn.Close()
}
