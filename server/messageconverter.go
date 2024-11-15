package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
)

// MessageResult is used to couple a possible error when sending a message
type MessageResult struct {
	sessionId string
	message   string
	err       error
}

// convertMessages converts incoming messages into their respective internal types
func convertMessages(ctx context.Context, incomingMessages <-chan MessageResult, textMessages chan<- TextMessage, commands chan<- Command) {
	for {
		select {
		case <-ctx.Done():
			return
		case incomingMessage := <-incomingMessages:
			message := cleanIncomingMessageString(incomingMessage.message)
			slog.Debug("incoming message", "message", message)
			if incomingMessage.err != nil {
				slog.Warn("incoming message error", "err", incomingMessage.err)
				if errors.Is(incomingMessage.err, io.EOF) {
					commands <- Command{sessionId: incomingMessage.sessionId, commandType: Quit, arguments: nil}
				}
				continue
			}
			if strings.HasPrefix(message, "/") {
				command := strings.TrimPrefix(message, "/")
				commandSplit := strings.Split(command, " ")
				commandType := commandSplit[0]
				commandArgs := commandSplit[1:]
				commands <- Command{sessionId: incomingMessage.sessionId, commandType: MatchCommandTypeStringToCommandType(commandType), arguments: commandArgs}
			} else {
				textMessages <- TextMessage{sessionId: incomingMessage.sessionId, message: message}
			}
		}
	}
}

// cleanIncomingMessageString is a helper function to clean strings that were received by the client
func cleanIncomingMessageString(message string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(message, "\n"), "\r"))
}
