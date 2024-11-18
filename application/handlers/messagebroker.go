package handlers

import (
	"context"

	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

// HandleMessages handles all incoming messages.
func HandleMessages(ctx context.Context, sessions <-chan domain.Session, textMessages <-chan domain.TextMessage, commands <-chan domain.Command) {
	sessionRepository := domain.NewInMemorySessionRepository()
	userRepository := domain.NewInMemoryUserRepository()
	userSessionRepository := domain.NewInMemoryUserSessionRepository()
	chatService := application.NewChatService(sessionRepository, userRepository, userSessionRepository)
	for {
		select {
		case <-ctx.Done():
			return
		case newSession := <-sessions:
			HandleNewSession(newSession, chatService)
		case textMessage := <-textMessages:
			HandleTextMessage(textMessage, chatService)
		case command := <-commands:
			HandleCommand(command, chatService)
		}
	}
}
