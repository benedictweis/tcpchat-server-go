package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"tcpchat-server-go/domain"
)

type TCPChatServer struct {
	address net.TCPAddr
}

// NewTCPChatServer creates a new instance of TCPChatServer with an address and a port
func NewTCPChatServer(address string, port int) (*TCPChatServer, error) {
	tcpAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}
	return &TCPChatServer{address: *tcpAddress}, nil
}

// Start starts the TCPChatServer instance and returns when ctx is Done
func (t *TCPChatServer) Start(ctx context.Context) error {
	listener, err := net.ListenTCP("tcp", &t.address)
	if err != nil {
		return err
	}
	defer listener.Close()
	messagesRead := make(chan MessageResult, 5) // Buffer to allow for bursts when sending messages
	sessions := make(chan domain.Session)
	textMessages := make(chan TextMessage)
	commands := make(chan Command)
	var activeConnections sync.WaitGroup
	go convertMessages(ctx, messagesRead, textMessages, commands)
	go handleMessages(ctx, sessions, textMessages, commands)
	go handleConnections(ctx, listener, &activeConnections, messagesRead, sessions)
	<-ctx.Done()
	log.Println("info: context is done, waiting for active connections to be closed")
	activeConnections.Wait()
	log.Println("info: active connections closed, stopping the server")
	return nil
}