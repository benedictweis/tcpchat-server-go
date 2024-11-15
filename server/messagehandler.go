package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

// TextMessage represents a message a user intends to send
type TextMessage struct {
	sessionId string
	message   string
}

type CommandType int

const (
	Unknown CommandType = iota
	ChangeName
	PrivateMessage
	CreateAccount
	Login
	ChangePassword
	Info
	Who
	Quit
)

// String implements the string variants of CommandType
func (c CommandType) String() string {
	return [...]string{"unknown", "name", "msg", "acc", "login", "passwd", "info", "who", "quit"}[c]
}

// MatchCommandTypeStringToCommandType is used to mach a received command as a string to the CommandType used to communicate the command
func MatchCommandTypeStringToCommandType(s string) CommandType {
	for currentCommandType := Unknown; currentCommandType <= Quit; currentCommandType++ {
		if currentCommandType.String() == s {
			return currentCommandType
		}
	}
	return Unknown
}

// Command represents a command a user wants to be executed
type Command struct {
	sessionId   string
	commandType CommandType
	arguments   []string
}

// handleMessages handles all incoming messages
func handleMessages(ctx context.Context, sessions <-chan domain.Session, textMessages <-chan TextMessage, commands <-chan Command) {
	sessionRepository := domain.NewInMemorySessionRepository()
	userRepository := domain.NewInMemoryUserRepository()
	userSessionRepository := domain.NewInMemoryUserSessionRepository()
	chatService := application.NewChatService(sessionRepository, userRepository, userSessionRepository)
	for {
		select {
		case <-ctx.Done():
			return
		case newSession := <-sessions:
			slog.Info("received new session", "sessionId", newSession.Id)
			chatService.RegisterNewSession(newSession)
			slog.Info("registered new session", "sessionId", newSession.Id)
			chatService.SendMessageToSessionFromServer(newSession.Id, "Welcome to this server!")
		case textMessage := <-textMessages:
			slog.Info("received text message", "sessionId", textMessage.sessionId, "textMessage", textMessage.message)
			err := chatService.SendTextMessageToEveryone(textMessage.sessionId, textMessage.message)
			if err != nil {
				handleErrors(err, chatService, textMessage.sessionId)
				continue
			}
			slog.Info("sent text message from to everyone", "sessionId", textMessage.sessionId, "textMessage", textMessage.message)
		case command := <-commands:
			slog.Info("received command", "sessionId", command.sessionId, "commandType", command.commandType, "commandArgs", command.arguments)
			switch command.commandType {
			case Unknown:
				slog.Info("received unknown command", "sessionId", command.sessionId, "commandType", command.commandType, "commandArgs", command.arguments)
				chatService.SendMessageToSessionFromServer(command.sessionId, "Unknown command")
			case ChangeName:
				if len(command.arguments) != 1 {
					slog.Info("invalid number of arguments", "sessionId", command.sessionId, "commandType", command.commandType, "commandArgs", command.arguments)
					chatService.SendMessageToSessionFromServer(command.sessionId, "Wrong number of arguments, usage: /name <new username>")
					continue
				}
				newUserName := command.arguments[0]
				err := chatService.ChangeUserName(command.sessionId, newUserName)
				if err != nil {
					handleErrors(err, chatService, command.sessionId)
					continue
				}
				slog.Info("changed name of user", "sessionId", command.sessionId, "newUserName", newUserName)
				chatService.SendMessageToSessionFromServer(command.sessionId, fmt.Sprintf("Changed username to %s", newUserName))
			case PrivateMessage:
				if len(command.arguments) < 2 {
					slog.Info("invalid number of arguments", "sessionId", command.sessionId, "commandType", command.commandType, "commandArgs", command.arguments)
					chatService.SendMessageToSessionFromServer(command.sessionId, "Wrong number of arguments, usage: /msg <username> <message...>")
					continue
				}
				messagePartnerUserName := command.arguments[0]
				message := strings.Join(command.arguments[1:], " ")
				err := chatService.SendPrivateMessage(command.sessionId, messagePartnerUserName, message)
				if err != nil {
					handleErrors(err, chatService, command.sessionId)
					continue
				}
				slog.Info("sent private message", "sessionId", command.sessionId, "messagePartnerUserName", messagePartnerUserName)
			case CreateAccount:
				if len(command.arguments) != 2 {
					slog.Info("invalid number of arguments", "sessionId", command.sessionId, "commandType", command.commandType, "commandArgs", command.arguments)
					chatService.SendMessageToSessionFromServer(command.sessionId, "Wrong number of arguments, usage: /acc <username> <password>")
					continue
				}
				userName := command.arguments[0]
				password := command.arguments[1]
				err := chatService.CreateAccount(command.sessionId, userName, password)
				if err != nil {
					handleErrors(err, chatService, command.sessionId)
					continue
				}
				slog.Info("created new account", "userName", userName)
				chatService.SendMessageToSessionFromServer(command.sessionId, "Created new account, please login now")
			case Login:
				if len(command.arguments) != 2 {
					slog.Info("invalid number of arguments", "sessionId", command.sessionId, "commandType", command.commandType, "commandArgs", command.arguments)
					chatService.SendMessageToSessionFromServer(command.sessionId, "Wrong number of arguments, usage: /login <username> <password>")
					continue
				}
				userName := command.arguments[0]
				password := command.arguments[1]

				err := chatService.Login(command.sessionId, userName, password)
				if err != nil {
					handleErrors(err, chatService, command.sessionId)
					continue
				}
				slog.Info("logged in session", "sessionId", command.sessionId, "userName", userName)
				chatService.SendMessageToSessionFromServer(command.sessionId, "Logged in")
			case ChangePassword:
				if len(command.arguments) != 2 {
					slog.Info("invalid number of arguments", "sessionId", command.sessionId, "commandType", command.commandType, "commandArgs", command.arguments)
					chatService.SendMessageToSessionFromServer(command.sessionId, "Wrong number of arguments, usage: /passwd <old password> <new password>")
					continue
				}
				oldPassword := command.arguments[0]
				newPassword := command.arguments[1]
				err := chatService.ChangePassword(command.sessionId, oldPassword, newPassword)
				if err != nil {
					handleErrors(err, chatService, command.sessionId)
					continue
				}
				slog.Info("changed password of user associated with session", "sessionId", command.sessionId)
				chatService.SendMessageToSessionFromServer(command.sessionId, "Changed Password")
			case Info:
				userName := chatService.GetUserNameForSessionId(command.sessionId)
				slog.Info("served info", "sessionId", command.sessionId)
				chatService.SendMessageToSessionFromServer(command.sessionId, fmt.Sprintf("sessionId: %s\n[server] userName:  %s\n", command.sessionId, userName))
			case Who:
				for _, userName := range chatService.GetAllLoggedInUserNames() {
					chatService.SendMessageToSessionFromServer(command.sessionId, userName)
				}
				slog.Info("served who", "sessionId", command.sessionId)
			case Quit:
				chatService.QuitSession(command.sessionId)
				slog.Info("quit session", "sessionId", command.sessionId)
			}
		}
	}
}

func handleErrors(err error, chatService *application.ChatService, sessionId string) {
	var userFriendlyError application.UserFriendlyError
	if errors.As(err, &userFriendlyError) {
		slog.Info("recovered from error", "err", err)
		chatService.SendMessageToSessionFromServer(sessionId, userFriendlyError.UserFriendlyError())
	} else {
		slog.Error("internal server error", "err", err)
		chatService.SendMessageToSessionFromServer(sessionId, "Internal server error")
	}
}
