// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package application

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"

	"tcpchat-server-go/domain"
)

// MessageResult is used to couple a possible error when sending a Message.
type MessageResult struct {
	SessionID string
	Message   string
	Err       error
}

// ConvertMessages converts incoming messages into their respective internal types.
func ConvertMessages(ctx context.Context, incomingMessages <-chan MessageResult, textMessages chan<- domain.TextMessage, commands chan<- domain.Command) {
	for {
		select {
		case <-ctx.Done():
			return
		case incomingMessage := <-incomingMessages:
			message := cleanIncomingMessageString(incomingMessage.Message)
			slog.Debug("incoming Message", "Message", message)
			if incomingMessage.Err != nil {
				slog.Warn("incoming Message error", "Err", incomingMessage.Err)
				if errors.Is(incomingMessage.Err, io.EOF) {
					commands <- domain.Command{SessionID: incomingMessage.SessionID, CommandType: domain.Quit, Arguments: nil}
				}
				continue
			}
			if strings.HasPrefix(message, "/") {
				command := strings.TrimPrefix(message, "/")
				commandSplit := strings.Fields(command)
				commandType := commandSplit[0]
				commandArgs := commandSplit[1:]
				commands <- domain.Command{SessionID: incomingMessage.SessionID, CommandType: domain.MatchCommandTypeStringToCommandType(commandType), Arguments: commandArgs}
			} else {
				textMessages <- domain.TextMessage{SessionID: incomingMessage.SessionID, Message: message}
			}
		}
	}
}

// cleanIncomingMessageString is a helper function to clean strings that were received by the client.
func cleanIncomingMessageString(message string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(message, "\n"), "\r"))
}
