package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type TCPConn struct {
	conn net.Conn
}

func newTCPConn(conn net.Conn) *TCPConn {
	return &TCPConn{
		conn: conn,
	}
}

func (tc *TCPConn) sendMessage(data []byte) error {
	length := uint32(len(data))
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, length)
	
	_, err := tc.conn.Write(lengthBytes)
	if err != nil {
		return fmt.Errorf("failed to write message length: %v", err)
	}
	
	_, err = tc.conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write message data: %v", err)
	}
	
	return nil
}

func (tc *TCPConn) receiveMessage() ([]byte, error) {
	lengthBytes := make([]byte, 4)
	_, err := tc.conn.Read(lengthBytes)

	length := binary.BigEndian.Uint32(lengthBytes)

	if length > 10*1024*1024 {
		caba_err("message too large")
		return nil, fmt.Errorf("failed to read message length: %v", err)
	}

	buffer := make([]byte, length)
	n, err := tc.conn.Read(buffer)
	if err != nil {
		caba_err("read tcp error")
		return nil, fmt.Errorf("failed to read: %v", err)
	}
	
	return buffer[:n], nil
}

func (tc *TCPConn) close_tcp() error {
	return tc.conn.Close()
}

func (tc *TCPConn) remoteAddr() net.Addr {
	return tc.conn.RemoteAddr()
}
