package server

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net"
	"sync"
	"tcpchat-server-go/domain"
)

// ConnectionResult is used to couple a possible error when accepting a connection with its result
type ConnectionResult struct {
	connection net.Conn
	err        error
}

// handleConnections is used to couple a possible error when accepting a connection with its result
func handleConnections(ctx context.Context, listener net.Listener, activeConnections *sync.WaitGroup, messagesRead chan<- MessageResult, sessions chan<- domain.Session) {
	connections := generateConnections(ctx, listener)
	for {
		select {
		case <-ctx.Done():
			return
		case connectionResult := <-connections:
			if connectionResult.err != nil {
				log.Printf("error: error accepting connection: %v", connectionResult.err)
				continue
			}
			go handleConnection(ctx, connectionResult.connection, activeConnections, sessions, messagesRead)
		}
	}
}

// generateConnections is used to accept incoming connections and send them on to a channel
func generateConnections(ctx context.Context, listener net.Listener) <-chan ConnectionResult {
	connections := make(chan ConnectionResult, 5) // Buffer of 5 to allow minor burst handling

	go func() {
		defer close(connections)
		for {
			conn, err := listener.Accept()
			select {
			case <-ctx.Done():
				return
			case connections <- ConnectionResult{connection: conn, err: err}:
			}
		}
	}()

	return connections
}

// handleConnection handles a single connection along with reading to and writing from the connection
func handleConnection(ctx context.Context, connection net.Conn, activeConnections *sync.WaitGroup, sessions chan<- domain.Session, readMessages chan<- MessageResult) {
	messagesToSession := make(chan string)
	sessionId := uuid.New().String()
	defer func() {
		log.Printf("info: closing connection with sessionId %s from %v", sessionId, connection.RemoteAddr())
		connection.Close()
		activeConnections.Done()
	}()
	activeConnections.Add(1)
	log.Printf("info: new connection established with sessionId %s from %v", sessionId, connection.RemoteAddr())
	closeSession := make(chan interface{})
	sessions <- domain.Session{sessionId, messagesToSession, closeSession}

	localCtx, closeLocalCtx := context.WithCancel(ctx)
	defer closeLocalCtx()

	go handleRead(localCtx, connection, readMessages, sessionId)
	go handleWrite(localCtx, connection, messagesToSession)

	select {
	case <-ctx.Done():
	case <-closeSession:
	}

	return
}
