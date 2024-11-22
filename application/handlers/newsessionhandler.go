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

func HandleNewSession(newSession domain.Session, chatService application.ChatService) {
	slog.Info("received new session", "sessionID", newSession.ID)
	chatService.RegisterNewSession(newSession)
	slog.Info("registered new session", "sessionID", newSession.ID)
	chatService.SendMessageToSessionFromServer(newSession.ID, "Welcome to this server!")
}
