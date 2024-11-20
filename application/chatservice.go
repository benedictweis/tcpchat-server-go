// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package application

import (
	"fmt"

	"github.com/benedictweis/tcpchat-server-go/domain"
)

//go:generate mockgen -destination=../test/mock/chatservice_mock.go . ChatService
type ChatService interface {
	SendMessageToSessionFromServer(sessionID string, message string)
	RegisterNewSession(newSession domain.Session)
	SendTextMessageToEveryone(sessionID, message string) error
	ChangeUserName(sessionID string, newUserName string) error
	SendPrivateMessage(sessionID, messagePartnerUserName, message string) error
	CreateAccount(sessionID, userName, password string) error
	Login(sessionID, userName, password string) error
	ChangePassword(sessionID, oldPassword, newPassword string) error
	GetUserNameForSessionID(sessionID string) string
	GetAllLoggedInUserNames() []string
	QuitSession(sessionID string)
}

type BasicChatService struct {
	sessionRepository     domain.SessionRepository
	userRepository        domain.UserRepository
	userSessionRepository domain.UserSessionRepository
}

func NewChatService(sessionRepository domain.SessionRepository, userRepository domain.UserRepository, userSessionRepository domain.UserSessionRepository) *BasicChatService {
	return &BasicChatService{sessionRepository: sessionRepository, userRepository: userRepository, userSessionRepository: userSessionRepository}
}

func (c BasicChatService) SendMessageToSessionFromServer(sessionID string, message string) {
	c.sendMessageToSession(sessionID, fmt.Sprintf("[plugin] %s", message))
}

func (c BasicChatService) sendMessageToSession(sessionID string, message string) {
	session, sessionExists := c.sessionRepository.FindByID(sessionID)
	if !sessionExists {
		return
	}
	session.MessagesToSession <- fmt.Sprintf("%s\n", message)
}

func (c BasicChatService) RegisterNewSession(newSession domain.Session) {
	c.sessionRepository.Add(newSession)
}

func (c BasicChatService) SendTextMessageToEveryone(sessionID, message string) error {
	_, sessionExists := c.sessionRepository.FindByID(sessionID)
	if !sessionExists {
		return fmt.Errorf("revieced a text Message from an unknown session id: %s", sessionID)
	}
	userSession, userSessionExists := c.userSessionRepository.FindBySessionID(sessionID)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionID)
	}
	userID := userSession.UserID
	user, userExists := c.userRepository.FindByID(userID)
	if !userExists {
		return fmt.Errorf("user was not found, userID: %s", userID)
	}
	otherSessions := c.sessionRepository.FindAllExceptBySessionID(sessionID)
	for _, otherSession := range otherSessions {
		c.sendMessageToSession(otherSession.ID, fmt.Sprintf("[%s] %s", user.Name, message))
	}
	return nil
}

func (c BasicChatService) ChangeUserName(sessionID string, newUserName string) error {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionID(sessionID)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionID)
	}
	user, userExists := c.userRepository.FindByID(userSession.UserID)
	if !userExists {
		return fmt.Errorf("user was not found, userID: %s", userSession.UserID)
	}
	user.Name = newUserName
	return nil
}

func (c BasicChatService) SendPrivateMessage(sessionID, messagePartnerUserName, message string) error {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionID(sessionID)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionID)
	}
	user, userExists := c.userRepository.FindByID(userSession.UserID)
	if !userExists {
		return fmt.Errorf("user was not found, userID: %s", userSession.UserID)
	}
	messagePartnerUser, messagePartnerUserExists := c.userRepository.FindByName(messagePartnerUserName)
	if !messagePartnerUserExists {
		return NewErrMessagePartnerDoesNotExist(sessionID, messagePartnerUserName)
	}
	messagePartnerUserSessions := c.userSessionRepository.FindByUserID(messagePartnerUser.ID)
	if len(messagePartnerUserSessions) == 0 {
		return NewErrMessagePartnerNotLoggedIn(sessionID, messagePartnerUserName)
	}
	for _, partnerUserSession := range messagePartnerUserSessions {
		c.sendMessageToSession(partnerUserSession.SessionID, fmt.Sprintf("[p %s] %s", user.Name, message))
	}
	return nil
}

func (c BasicChatService) CreateAccount(sessionID, userName, password string) error {
	user, err := domain.NewUser(userName, password)
	if err != nil {
		return NewErrCouldNotCreateUser(sessionID)
	}
	addedUser := c.userRepository.Add(user)
	if !addedUser {
		return NewErrUserNameAlreadyExists(sessionID, userName)
	}
	return nil
}

func (c BasicChatService) Login(sessionID, userName, password string) error {
	user, userExists := c.userRepository.FindByName(userName)
	if !userExists {
		return NewErrUserDoesNotExist(sessionID, userName)
	}
	passwordIsValid := user.PasswordIsValid(password)
	if !passwordIsValid {
		return NewErrPasswordIsInvalid(sessionID)
	}
	userSession := domain.NewUserSession(user.ID, sessionID)
	c.userSessionRepository.Add(userSession)
	return nil
}

func (c BasicChatService) ChangePassword(sessionID, oldPassword, newPassword string) error {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionID(sessionID)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionID)
	}
	user, userExists := c.userRepository.FindByID(userSession.UserID)
	if !userExists {
		return fmt.Errorf("user was not found, userID: %s", userSession.UserID)
	}
	if !user.PasswordIsValid(oldPassword) {
		return NewErrPasswordIsInvalid(sessionID)
	}
	err := user.SetPassword(newPassword)
	if err != nil {
		return NewErrPasswordIsInvalid(sessionID)
	}
	return nil
}

func (c BasicChatService) GetUserNameForSessionID(sessionID string) string {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionID(sessionID)
	if !userSessionExists {
		return ""
	}
	user, userExists := c.userRepository.FindByID(userSession.UserID)
	if !userExists {
		return ""
	}
	return user.Name
}

func (c BasicChatService) GetAllLoggedInUserNames() []string {
	userNames := make([]string, 0)
	for _, user := range c.userRepository.GetAll() {
		userSessions := c.userSessionRepository.DeleteByUserID(user.ID)
		if len(userSessions) > 0 {
			userNames = append(userNames, user.Name)
		}
	}
	return userNames
}

func (c BasicChatService) QuitSession(sessionID string) {
	session, sessionExists := c.sessionRepository.FindByID(sessionID)
	if !sessionExists {
		return
	}
	session.Close <- struct{}{}
	c.userSessionRepository.DeleteBySessionID(sessionID)
	c.sessionRepository.Delete(sessionID)
}
