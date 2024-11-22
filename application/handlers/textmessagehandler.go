// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package handlers

import (
	"log/slog"

	"github.com/benedictweis/tcpchat-server-go/application"
	"github.com/benedictweis/tcpchat-server-go/domain"
)

func HandleTextMessage(textMessage domain.TextMessage, chatService application.ChatService) {
	slog.Info("received text message", "sessionID", textMessage.SessionID, "textMessage", textMessage.Message)
	err := chatService.SendTextMessageToEveryone(textMessage.SessionID, textMessage.Message)
	if err != nil {
		handleErrors(err, chatService, textMessage.SessionID)
		return
	}
	slog.Info("sent text message from to everyone", "sessionID", textMessage.SessionID, "textMessage", textMessage.Message)
}
