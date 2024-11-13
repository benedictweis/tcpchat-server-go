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
	Info
	Quit
)

// String implements the string variants of CommandType
func (c CommandType) String() string {
	return [...]string{"unknown", "name", "msg", "account", "login", "password", "info", "quit"}[c]
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
				userName := command.arguments[0]
				password := command.arguments[1]
				user, err := domain.NewUser(userName, password)
				if err != nil {
					log.Printf("warn: failed to create user: %v\n", err)
					session.MessagesToSession <- "[server] Failed to create user!\n"
					continue
				}
				userExists := userRepository.Add(user)
				if !userExists {
					log.Printf("info: failed to create user %s, user already exists\n", userName)
					session.MessagesToSession <- "[server] Failed to create user, user already exists!\n"
					continue
				}
				session.MessagesToSession <- "[server] Created new account, please login now!\n"
			case Login:
				userName := command.arguments[0]
				password := command.arguments[1]
				user, userExists := userRepository.FindByName(userName)
				if !userExists {
					log.Printf("info: failed to find user by name: %s\n", userName)
					session.MessagesToSession <- "[server] User does not exist!\n"
					continue
				}
				passwordIsValid := user.PasswordIsValid(password)
				if !passwordIsValid {
					log.Printf("info: invalid password for user %s\n", userName)
					session.MessagesToSession <- "[server] Invalid password!\n"
					continue
				}
				userSession := domain.NewUserSession(command.sessionId, userName)
				userSessionRepository.Add(userSession)
				log.Printf("info: logged in user: %s\n", userName)
				session.MessagesToSession <- "[server] Logged in!\n"
			case ChangePassword:
				// TODO
				session.MessagesToSession <- "[server] Unimplemented!\n"
			case Info:
				userSession, userSessionExists := userSessionRepository.FindBySessionId(command.sessionId)
				if !userSessionExists {
					log.Printf("info: on retrieving user info: not logged in with session: %s\n", command.sessionId)
					session.MessagesToSession <- "[server] User not logged in!\n"
					continue
				}
				session.MessagesToSession <- fmt.Sprintf("[server] sessionId: %s\n[server] userName:  %s\n", userSession.SessionId, userSession.UserName)
			case Quit:
				session.Close <- struct{}{}
				sessionRepository.Delete(command.sessionId)
			}
		}
	}
}
