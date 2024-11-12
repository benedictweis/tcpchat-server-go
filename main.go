package main

import (
	"context"
	"log"
	"tcpchat-server-go/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tcpChatServer, err := server.NewTCPChatServer("localhost", 8080)
	if err != nil {
		log.Printf("error: failed to initialize tcp chat server: %v", err)
		return
	}
	err = tcpChatServer.Start(ctx)
	if err != nil {
		log.Printf("error: failed to start tcp chat server: %v", err)
		return
	}
}
