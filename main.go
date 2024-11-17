package main

import (
	"context"
	"log/slog"
	"os"
	"tcpchat-server-go/plugin"
)

func main() {
	setupLogging()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tcpChatServer, err := plugin.NewTCPChatServer("localhost", 8080)
	if err != nil {
		slog.Error("failed to initialize tcp chat plugin", "err", err)
		return
	}
	err = tcpChatServer.Start(ctx)
	if err != nil {
		slog.Error("error: failed to start tcp chat plugin", "err", err)
		return
	}
}

func setupLogging() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
}
