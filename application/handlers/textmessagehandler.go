package handlers

import (
	"log/slog"

	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

func HandleTextMessage(textMessage domain.TextMessage, chatService *application.ChatService) {
	slog.Info("received text message", "sessionID", textMessage.SessionID, "textMessage", textMessage.Message)
	err := chatService.SendTextMessageToEveryone(textMessage.SessionID, textMessage.Message)
	if err != nil {
		handleErrors(err, chatService, textMessage.SessionID)
		return
	}
	slog.Info("sent text message from to everyone", "sessionID", textMessage.SessionID, "textMessage", textMessage.Message)
}
