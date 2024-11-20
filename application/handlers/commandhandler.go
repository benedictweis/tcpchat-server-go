// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package handlers

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/benedictweis/tcpchat-server-go/application"
	"github.com/benedictweis/tcpchat-server-go/domain"
)

func HandleCommand(command domain.Command, chatService *application.BasicChatService) {
	slog.Info("received command", "sessionID", command.SessionID, "commandType", command.CommandType, "commandArgs", command.Arguments)
	matchCommandTypeToFunc(command.CommandType)(command, chatService)
}

func matchCommandTypeToFunc(commandType domain.CommandType) func(command domain.Command, chatService *application.BasicChatService) {
	handlers := []func(command domain.Command, chatService *application.BasicChatService){
		handleUnknownCommand,        // 0
		handleChangeNameCommand,     // 1
		handlePrivateMessageCommand, // 2
		handleCreateAccountCommand,  // 3
		handleLoginCommand,          // 4
		handleChangePasswordCommand, // 5
		handleInfoCommand,           // 6
		handleWhoCommand,            // 7
		handleQuitCommand,           // 8
	}

	// Ensure commandType is valid and within bounds
	if int(commandType) < 0 || int(commandType) >= len(handlers) {
		return handleUnknownCommand
	}

	return handlers[commandType]
}

func handleUnknownCommand(command domain.Command, chatService *application.BasicChatService) {
	slog.Info("received unknown command", "sessionID", command.SessionID, "commandType", command.CommandType, "commandArgs", command.Arguments)
	chatService.SendMessageToSessionFromServer(command.SessionID, "Unknown command")
}

func handleChangeNameCommand(command domain.Command, chatService *application.BasicChatService) {
	if len(command.Arguments) != 1 {
		slog.Info("invalid number of arguments", "sessionID", command.SessionID, "commandType", command.CommandType, "commandArgs", command.Arguments)
		chatService.SendMessageToSessionFromServer(command.SessionID, "Wrong number of arguments, usage: /name <new username>")
		return
	}
	newUserName := command.Arguments[0]
	err := chatService.ChangeUserName(command.SessionID, newUserName)
	if err != nil {
		handleErrors(err, chatService, command.SessionID)
		return
	}
	slog.Info("changed name of user", "sessionID", command.SessionID, "newUserName", newUserName)
	chatService.SendMessageToSessionFromServer(command.SessionID, fmt.Sprintf("Changed username to %s", newUserName))
}

func handlePrivateMessageCommand(command domain.Command, chatService *application.BasicChatService) {
	if len(command.Arguments) < 2 {
		slog.Info("invalid number of arguments", "sessionID", command.SessionID, "commandType", command.CommandType, "commandArgs", command.Arguments)
		chatService.SendMessageToSessionFromServer(command.SessionID, "Wrong number of arguments, usage: /msg <username> <message...>")
		return
	}
	messagePartnerUserName := command.Arguments[0]
	message := strings.Join(command.Arguments[1:], " ")
	err := chatService.SendPrivateMessage(command.SessionID, messagePartnerUserName, message)
	if err != nil {
		handleErrors(err, chatService, command.SessionID)
		return
	}
	slog.Info("sent private message", "sessionID", command.SessionID, "messagePartnerUserName", messagePartnerUserName)
}

func handleCreateAccountCommand(command domain.Command, chatService *application.BasicChatService) {
	if len(command.Arguments) != 2 {
		slog.Info("invalid number of arguments", "sessionID", command.SessionID, "commandType", command.CommandType, "commandArgs", command.Arguments)
		chatService.SendMessageToSessionFromServer(command.SessionID, "Wrong number of arguments, usage: /acc <username> <password>")
		return
	}
	userName := command.Arguments[0]
	password := command.Arguments[1]
	err := chatService.CreateAccount(command.SessionID, userName, password)
	if err != nil {
		handleErrors(err, chatService, command.SessionID)
		return
	}
	slog.Info("created new account", "userName", userName)
	chatService.SendMessageToSessionFromServer(command.SessionID, "Created new account, please login now")
}

func handleLoginCommand(command domain.Command, chatService *application.BasicChatService) {
	if len(command.Arguments) != 2 {
		slog.Info("invalid number of arguments", "sessionID", command.SessionID, "commandType", command.CommandType, "commandArgs", command.Arguments)
		chatService.SendMessageToSessionFromServer(command.SessionID, "Wrong number of arguments, usage: /login <username> <password>")
		return
	}
	userName := command.Arguments[0]
	password := command.Arguments[1]

	err := chatService.Login(command.SessionID, userName, password)
	if err != nil {
		handleErrors(err, chatService, command.SessionID)
		return
	}
	slog.Info("logged in session", "sessionID", command.SessionID, "userName", userName)
	chatService.SendMessageToSessionFromServer(command.SessionID, "Logged in")
}

func handleChangePasswordCommand(command domain.Command, chatService *application.BasicChatService) {
	if len(command.Arguments) != 2 {
		slog.Info("invalid number of arguments", "sessionID", command.SessionID, "commandType", command.CommandType, "commandArgs", command.Arguments)
		chatService.SendMessageToSessionFromServer(command.SessionID, "Wrong number of arguments, usage: /passwd <old password> <new password>")
		return
	}
	oldPassword := command.Arguments[0]
	newPassword := command.Arguments[1]
	err := chatService.ChangePassword(command.SessionID, oldPassword, newPassword)
	if err != nil {
		handleErrors(err, chatService, command.SessionID)
		return
	}
	slog.Info("changed password of user associated with session", "sessionID", command.SessionID)
	chatService.SendMessageToSessionFromServer(command.SessionID, "Changed Password")
}

func handleInfoCommand(command domain.Command, chatService *application.BasicChatService) {
	userName := chatService.GetUserNameForSessionID(command.SessionID)
	slog.Info("served info", "sessionID", command.SessionID)
	chatService.SendMessageToSessionFromServer(command.SessionID, fmt.Sprintf("sessionID: %s\n[plugin] userName:  %s", command.SessionID, userName))
}

func handleWhoCommand(command domain.Command, chatService *application.BasicChatService) {
	for _, userName := range chatService.GetAllLoggedInUserNames() {
		chatService.SendMessageToSessionFromServer(command.SessionID, userName)
	}
	slog.Info("served who", "sessionID", command.SessionID)
}

func handleQuitCommand(command domain.Command, chatService *application.BasicChatService) {
	chatService.QuitSession(command.SessionID)
	slog.Info("quit session", "sessionID", command.SessionID)
}
