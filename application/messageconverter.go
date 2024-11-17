package application

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"tcpchat-server-go/domain"
)

// MessageResult is used to couple a possible error when sending a Message
type MessageResult struct {
	SessionId string
	Message   string
	Err       error
}

// ConvertMessages converts incoming messages into their respective internal types
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
					commands <- domain.Command{SessionId: incomingMessage.SessionId, CommandType: domain.Quit, Arguments: nil}
				}
				continue
			}
			if strings.HasPrefix(message, "/") {
				command := strings.TrimPrefix(message, "/")
				commandSplit := strings.Split(command, " ")
				commandType := commandSplit[0]
				commandArgs := commandSplit[1:]
				commands <- domain.Command{SessionId: incomingMessage.SessionId, CommandType: domain.MatchCommandTypeStringToCommandType(commandType), Arguments: commandArgs}
			} else {
				textMessages <- domain.TextMessage{SessionId: incomingMessage.SessionId, Message: message}
			}
		}
	}
}

// cleanIncomingMessageString is a helper function to clean strings that were received by the client
func cleanIncomingMessageString(message string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(message, "\n"), "\r"))
}
