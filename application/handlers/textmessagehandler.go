package handlers

import (
	"log/slog"
	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

func HandleTextMessage(textMessage domain.TextMessage, chatService *application.ChatService) {
	slog.Info("received text message", "sessionId", textMessage.SessionId, "textMessage", textMessage.Message)
	err := chatService.SendTextMessageToEveryone(textMessage.SessionId, textMessage.Message)
	if err != nil {
		handleErrors(err, chatService, textMessage.SessionId)
		return
	}
	slog.Info("sent text message from to everyone", "sessionId", textMessage.SessionId, "textMessage", textMessage.Message)
}
