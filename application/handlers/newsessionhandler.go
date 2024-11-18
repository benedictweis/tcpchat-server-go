package handlers

import (
	"log/slog"

	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

func HandleNewSession(newSession domain.Session, chatService *application.ChatService) {
	slog.Info("received new session", "sessionID", newSession.ID)
	chatService.RegisterNewSession(newSession)
	slog.Info("registered new session", "sessionID", newSession.ID)
	chatService.SendMessageToSessionFromServer(newSession.ID, "Welcome to this plugin!")
}
