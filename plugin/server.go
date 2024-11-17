package plugin

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"tcpchat-server-go/application"
	"tcpchat-server-go/application/handlers"
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
	slog.Info("starting tcp chat plugin", "address", t.address.String())
	listener, err := net.ListenTCP("tcp", &t.address)
	if err != nil {
		return err
	}
	defer listener.Close()
	messagesRead := make(chan application.MessageResult, 5) // Buffer to allow for bursts when sending messages
	sessions := make(chan domain.Session)
	textMessages := make(chan domain.TextMessage)
	commands := make(chan domain.Command)
	var activeConnections sync.WaitGroup
	go application.ConvertMessages(ctx, messagesRead, textMessages, commands)
	go handlers.HandleMessages(ctx, sessions, textMessages, commands)
	go handleConnections(ctx, listener, &activeConnections, messagesRead, sessions)
	slog.Info("tcp chat is up", "address", t.address.String())
	<-ctx.Done()
	slog.Info("context is done, waiting for active connections to be closed", "address", t.address.String())
	activeConnections.Wait()
	slog.Info("active connections closed, stopping the plugin", "address", t.address.String())
	return nil
}
