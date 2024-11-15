package main

import (
	"context"
	"log/slog"
	"os"
	"tcpchat-server-go/server"
)

func main() {
	setupLogging()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tcpChatServer, err := server.NewTCPChatServer("localhost", 8080)
	if err != nil {
		slog.Error("failed to initialize tcp chat server", "err", err)
		return
	}
	err = tcpChatServer.Start(ctx)
	if err != nil {
		slog.Error("error: failed to start tcp chat server", "err", err)
		return
	}
}

func setupLogging() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
