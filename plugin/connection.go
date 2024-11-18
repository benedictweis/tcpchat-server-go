package plugin

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

// ConnectionResult is used to couple a possible error when accepting a connection with its result.
type ConnectionResult struct {
	connection net.Conn
	err        error
}

// handleConnections is used to couple a possible error when accepting a connection with its result.
func handleConnections(ctx context.Context, listener net.Listener, activeConnections *sync.WaitGroup, messagesRead chan<- application.MessageResult, sessions chan<- domain.Session) {
	connections := generateConnections(ctx, listener)
	for {
		select {
		case <-ctx.Done():
			return
		case connectionResult := <-connections:
			if connectionResult.err != nil {
				slog.Error("error accepting connection", "err", connectionResult.err)
				continue
			}
			go handleConnection(ctx, connectionResult.connection, activeConnections, sessions, messagesRead)
		}
	}
}

// generateConnections is used to accept incoming connections and send them on to a channel.
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

// handleConnection handles a single connection along with reading to and writing from the connection.
func handleConnection(ctx context.Context, connection net.Conn, activeConnections *sync.WaitGroup, sessions chan<- domain.Session, readMessages chan<- application.MessageResult) {
	messagesToSession := make(chan string)
	closeSession := make(chan interface{})
	session := domain.NewSession(messagesToSession, closeSession)
	slog.Info("new connection established", "sessionID", session.ID, "remoteAddr", connection.RemoteAddr())
	defer func() {
		slog.Info("closing session", "sessionID", session.ID, "remoteAddr", connection.RemoteAddr())
		connection.Close()
		activeConnections.Done()
	}()
	activeConnections.Add(1)

	sessions <- *session

	localCtx, closeLocalCtx := context.WithCancel(ctx)
	defer closeLocalCtx()

	go handleRead(localCtx, connection, readMessages, session.ID)
	go handleWrite(localCtx, connection, messagesToSession)

	select {
	case <-ctx.Done():
	case <-closeSession:
	}
}
