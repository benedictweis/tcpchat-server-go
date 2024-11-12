package server

import (
	"context"
	"fmt"
	"log"
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
	Quit
)

// String implements the string variants of CommandType
func (c CommandType) String() string {
	return [...]string{"unknown", "name", "msg", "account", "login", "password", "quit"}[c]
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
	_ = domain.NewInMemoryUserRepository()
	for {
		select {
		case <-ctx.Done():
			return
		case newSession := <-sessions:
			log.Printf("info: recieved new session with sessionId %s\n", newSession.Id)
			newSession.MessagesToSession <- fmt.Sprintf("[server] Welcome to this server!\n")
			sessionRepository.Add(newSession)
		case textMessage := <-textMessages:
			log.Printf("info: recieved text message from %s, message: %s\n", textMessage.sessionId, textMessage.message)
		case command := <-commands:
			log.Printf("info: recieved command from %s, type: %v, arguments: %v\n", command.sessionId, command.commandType, command.arguments)
			session, _ := sessionRepository.FindById(command.sessionId)
			switch command.commandType {
			case Unknown:
				session.MessagesToSession <- "[server] Entered unknown command!\n"
			case ChangeName:
				// TODO
				session.MessagesToSession <- "[server] Unimplemented!\n"
			case PrivateMessage:
				// TODO
				session.MessagesToSession <- "[server] Unimplemented!\n"
			case CreateAccount:
				// TODO
				session.MessagesToSession <- "[server] Unimplemented!\n"
			case Login:
				// TODO
				session.MessagesToSession <- "[server] Unimplemented!\n"
			case ChangePassword:
				// TODO
				session.MessagesToSession <- "[server] Unimplemented!\n"
			case Quit:
				session.Close <- struct{}{}
				sessionRepository.Delete(command.sessionId)
			}
		}
	}
}
