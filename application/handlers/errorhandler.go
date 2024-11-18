package handlers

import (
	"errors"
	"log/slog"

	"tcpchat-server-go/application"
)

func handleErrors(err error, chatService *application.ChatService, sessionId string) {
	var userFriendlyError application.UserFriendlyError
	if errors.As(err, &userFriendlyError) {
		slog.Info("recovered from error", "err", err)
		chatService.SendMessageToSessionFromServer(sessionId, userFriendlyError.UserFriendlyError())
	} else {
		slog.Error("internal plugin error", "err", err)
		chatService.SendMessageToSessionFromServer(sessionId, "Internal plugin error")
	}
}
