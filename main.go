package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var maindb db = db{}
var cache_ cache = cache{m: map[string]string{}}
var q *queue = newQueue(16384)

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func main() {
	cfgerr := load_cfg("config")

	if cfgerr != nil {
		fmt.Println("config file not found. creating config.")
		file, err := os.Create("config")
		file.WriteString("PASSKEY=\"YbQjLuBXX4yv18oqmEzXOnf67USJJZN8\"\nCACHE_SIZE=\"8192\"")
		_check(err)
		defer file.Close()
	}

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
