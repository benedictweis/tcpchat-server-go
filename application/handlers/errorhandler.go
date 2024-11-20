// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package handlers

import (
	"errors"
	"log/slog"

	"github.com/benedictweis/tcpchat-server-go/application"
)

func handleErrors(err error, chatService application.ChatService, sessionID string) {
	var userFriendlyError application.UserFriendlyError
	if errors.As(err, &userFriendlyError) {
		slog.Info("recovered from error", "err", err)
		chatService.SendMessageToSessionFromServer(sessionID, userFriendlyError.UserFriendlyError())
	} else {
		slog.Error("internal server error", "err", err)
		chatService.SendMessageToSessionFromServer(sessionID, "internal server error")
	}
}
