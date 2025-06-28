package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var maindb db = db{}
var cache_ cache = cache{m: map[string]string{}}

const (
	HOST = "localhost"
	PORT = "8080"
	TYPE = "tcp"
)

func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + 
		   string(rune(time.Now().Nanosecond()%26+65))
}

type server struct {
	clients         map[string]*TCPConn
    clientsMutex    sync.RWMutex
	requestQueue    chan transaction
    listener        net.Listener
}

func (s *server) start(address string) error {
	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		return err
	}

	caba_log("TCP server listening")

	go s.processRequests()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			caba_err("Error accepting connection")
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
	tcpConn := NewTCPConn(conn)
	defer tcpConn.Close()

	clientID := generateClientID()
	s.addClient(clientID, tcpConn)
	defer s.removeClient(clientID)

	caba_log("Client connected")

	for {
		msg, err := tcpConn.receiveMessage()
		if err != nil {
			caba_err(err)
			break
		}
		reqs := parseRequests(string(msg))

		for _, v := range reqs {
			v.CID = clientID
			tr := transaction{
				-1,
				v,
			}

			select {
			case s.requestQueue <- tr:
			default:
				s.sendResponse(clientID, "Request queue is full")
			}
		}
	}
}

func (s *server) addClient(clientID string, conn *TCPConn) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	s.clients[clientID] = conn
}

func (s *server) removeClient(clientID string) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	if conn, exists := s.clients[clientID]; exists {
		conn.Close()
		delete(s.clients, clientID)
	}
}

func (s *server) sendResponse(clientID string, response string) {
	s.clientsMutex.RLock()
	conn, exists := s.clients[clientID]
	s.clientsMutex.RUnlock()

	if !exists {
		caba_err("Client not found")
		return
	}

	err := conn.SendMessage([]byte(response))
	if err != nil {
		caba_err("Error sending response to client")
		s.removeClient(clientID)
	}
}

func (s *server) stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func new_server() *server {
    return &server{
        clients:      make(map[string]*TCPConn),
        requestQueue: make(chan transaction, 1000),
    }
}

func main() {
	server := new_server()

	cfgerr := load_cfg("config")

	if cfgerr != nil {
		fmt.Println("config file not found. creating config.")
		file, err := os.Create("config")
		file.WriteString("PASSKEY=\"YbQjLuBXX4yv18oqmEzXOnf67USJJZN8\"\nCACHE_SIZE=\"8192\"")
		_check(err)
		defer file.Close()
	}

	address := ":8080"
	
	err := server.start(address)
	if err != nil {
		caba_err("Server failed to start")
	}
}

func (s *server) processRequests() {
	for tr := range s.requestQueue {
		text := tr.execute()
		
		s.sendResponse(tr.request.CID, text)
	}
}