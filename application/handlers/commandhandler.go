package handlers

import (
	"fmt"
	"log/slog"
	"strings"
	"tcpchat-server-go/application"
	"tcpchat-server-go/domain"
)

func HandleCommand(command domain.Command, chatService *application.ChatService) {
	slog.Info("received command", "sessionId", command.SessionId, "commandType", command.CommandType, "commandArgs", command.Arguments)
	switch command.CommandType {
	case domain.Unknown:
		slog.Info("received unknown command", "sessionId", command.SessionId, "commandType", command.CommandType, "commandArgs", command.Arguments)
		chatService.SendMessageToSessionFromServer(command.SessionId, "Unknown command")
	case domain.ChangeName:
		if len(command.Arguments) != 1 {
			slog.Info("invalid number of arguments", "sessionId", command.SessionId, "commandType", command.CommandType, "commandArgs", command.Arguments)
			chatService.SendMessageToSessionFromServer(command.SessionId, "Wrong number of arguments, usage: /name <new username>")
			return
		}
		newUserName := command.Arguments[0]
		err := chatService.ChangeUserName(command.SessionId, newUserName)
		if err != nil {
			handleErrors(err, chatService, command.SessionId)
			return
		}
		slog.Info("changed name of user", "sessionId", command.SessionId, "newUserName", newUserName)
		chatService.SendMessageToSessionFromServer(command.SessionId, fmt.Sprintf("Changed username to %s", newUserName))
	case domain.PrivateMessage:
		if len(command.Arguments) < 2 {
			slog.Info("invalid number of arguments", "sessionId", command.SessionId, "commandType", command.CommandType, "commandArgs", command.Arguments)
			chatService.SendMessageToSessionFromServer(command.SessionId, "Wrong number of arguments, usage: /msg <username> <message...>")
			return
		}
		messagePartnerUserName := command.Arguments[0]
		message := strings.Join(command.Arguments[1:], " ")
		err := chatService.SendPrivateMessage(command.SessionId, messagePartnerUserName, message)
		if err != nil {
			handleErrors(err, chatService, command.SessionId)
			return
		}
		slog.Info("sent private message", "sessionId", command.SessionId, "messagePartnerUserName", messagePartnerUserName)
	case domain.CreateAccount:
		if len(command.Arguments) != 2 {
			slog.Info("invalid number of arguments", "sessionId", command.SessionId, "commandType", command.CommandType, "commandArgs", command.Arguments)
			chatService.SendMessageToSessionFromServer(command.SessionId, "Wrong number of arguments, usage: /acc <username> <password>")
			return
		}
		userName := command.Arguments[0]
		password := command.Arguments[1]
		err := chatService.CreateAccount(command.SessionId, userName, password)
		if err != nil {
			handleErrors(err, chatService, command.SessionId)
			return
		}
		slog.Info("created new account", "userName", userName)
		chatService.SendMessageToSessionFromServer(command.SessionId, "Created new account, please login now")
	case domain.Login:
		if len(command.Arguments) != 2 {
			slog.Info("invalid number of arguments", "sessionId", command.SessionId, "commandType", command.CommandType, "commandArgs", command.Arguments)
			chatService.SendMessageToSessionFromServer(command.SessionId, "Wrong number of arguments, usage: /login <username> <password>")
			return
		}
		userName := command.Arguments[0]
		password := command.Arguments[1]

		err := chatService.Login(command.SessionId, userName, password)
		if err != nil {
			handleErrors(err, chatService, command.SessionId)
			return
		}
		slog.Info("logged in session", "sessionId", command.SessionId, "userName", userName)
		chatService.SendMessageToSessionFromServer(command.SessionId, "Logged in")
	case domain.ChangePassword:
		if len(command.Arguments) != 2 {
			slog.Info("invalid number of arguments", "sessionId", command.SessionId, "commandType", command.CommandType, "commandArgs", command.Arguments)
			chatService.SendMessageToSessionFromServer(command.SessionId, "Wrong number of arguments, usage: /passwd <old password> <new password>")
			return
		}
		oldPassword := command.Arguments[0]
		newPassword := command.Arguments[1]
		err := chatService.ChangePassword(command.SessionId, oldPassword, newPassword)
		if err != nil {
			handleErrors(err, chatService, command.SessionId)
			return
		}
		slog.Info("changed password of user associated with session", "sessionId", command.SessionId)
		chatService.SendMessageToSessionFromServer(command.SessionId, "Changed Password")
	case domain.Info:
		userName := chatService.GetUserNameForSessionId(command.SessionId)
		slog.Info("served info", "sessionId", command.SessionId)
		chatService.SendMessageToSessionFromServer(command.SessionId, fmt.Sprintf("sessionId: %s\n[plugin] userName:  %s\n", command.SessionId, userName))
	case domain.Who:
		for _, userName := range chatService.GetAllLoggedInUserNames() {
			chatService.SendMessageToSessionFromServer(command.SessionId, userName)
		}
		slog.Info("served who", "sessionId", command.SessionId)
	case domain.Quit:
		chatService.QuitSession(command.SessionId)
		slog.Info("quit session", "sessionId", command.SessionId)
	}
}
