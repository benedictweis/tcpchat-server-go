package handlers

import (
	"log/slog"

	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

func HandleNewSession(newSession domain.Session, chatService *application.ChatService) {
	slog.Info("received new session", "sessionId", newSession.Id)
	chatService.RegisterNewSession(newSession)
	slog.Info("registered new session", "sessionId", newSession.Id)
	chatService.SendMessageToSessionFromServer(newSession.Id, "Welcome to this plugin!")
}
