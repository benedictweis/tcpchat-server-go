package application

import (
	"errors"
	"fmt"
	"tcpchat-server-go/domain"
)

type ChatService struct {
	sessionRepository     domain.SessionRepository
	userRepository        domain.UserRepository
	userSessionRepository domain.UserSessionRepository
}

func NewChatService(sessionRepository domain.SessionRepository, userRepository domain.UserRepository, userSessionRepository domain.UserSessionRepository) *ChatService {
	return &ChatService{sessionRepository: sessionRepository, userRepository: userRepository, userSessionRepository: userSessionRepository}
}

func (c ChatService) SendMessageToSessionFromServer(sessionId string, message string) {
	c.sendMessageToSession(sessionId, fmt.Sprintf("[server] %s", message))
}

func (c ChatService) sendMessageToSession(sessionId string, message string) {
	session, sessionExists := c.sessionRepository.FindById(sessionId)
	if !sessionExists {
		return
	}
	session.MessagesToSession <- fmt.Sprintf("%s\n", message)
}

func (c ChatService) RegisterNewSession(newSession domain.Session) {
	c.sessionRepository.Add(newSession)
}

func (c ChatService) SendTextMessageToEveryone(sessionId, message string) error {
	_, sessionExists := c.sessionRepository.FindById(sessionId)
	if !sessionExists {
		return errors.New(fmt.Sprintf("revieced a text message from an unknown session id: %s", sessionId))
	}
	userSession, userSessionExists := c.userSessionRepository.FindBySessionId(sessionId)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionId)
	}
	userId := userSession.UserId
	user, userExists := c.userRepository.FindById(userId)
	if !userExists {
		return errors.New(fmt.Sprintf("user was not found, userId: %s", userId))
	}
	otherSessions := c.sessionRepository.FindAllExceptBySessionId(sessionId)
	for _, otherSession := range otherSessions {
		c.sendMessageToSession(otherSession.Id, fmt.Sprintf("[%s] %s", user.Name, message))
	}
	return nil
}

func (c ChatService) ChangeUserName(sessionId string, newUserName string) error {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionId(sessionId)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionId)
	}
	user, userExists := c.userRepository.FindById(userSession.UserId)
	if !userExists {
		return errors.New(fmt.Sprintf("user was not found, userId: %s", userSession.UserId))
	}
	user.Name = newUserName
	return nil
}

func (c ChatService) SendPrivateMessage(sessionId, messagePartnerUserName, message string) error {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionId(sessionId)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionId)
	}
	user, userExists := c.userRepository.FindById(userSession.UserId)
	if !userExists {
		return errors.New(fmt.Sprintf("user was not found, userId: %s", userSession.UserId))
	}
	messagePartnerUser, messagePartnerUserExists := c.userRepository.FindByName(messagePartnerUserName)
	if !messagePartnerUserExists {
		return NewErrMessagePartnerDoesNotExist(sessionId, messagePartnerUserName)
	}
	messagePartnerUserSessions := c.userSessionRepository.FindByUserId(messagePartnerUser.Id)
	if len(messagePartnerUserSessions) == 0 {
		return NewErrMessagePartnerNotLoggedIn(sessionId, messagePartnerUserName)
	}
	for _, partnerUserSession := range messagePartnerUserSessions {
		c.sendMessageToSession(partnerUserSession.SessionId, fmt.Sprintf("[p %s] %s", user.Name, message))
	}
	return nil
}

func (c ChatService) CreateAccount(sessionId, userName, password string) error {
	user, err := domain.NewUser(userName, password)
	if err != nil {
		return NewErrCouldNotCreateUser(sessionId)
	}
	addedUser := c.userRepository.Add(user)
	if !addedUser {
		return NewErrUserNameAlreadyExists(sessionId, userName)
	}
	return nil
}

func (c ChatService) Login(sessionId, userName, password string) error {
	user, userExists := c.userRepository.FindByName(userName)
	if !userExists {
		return NewErrUserDoesNotExist(sessionId, userName)
	}
	passwordIsValid := user.PasswordIsValid(password)
	if !passwordIsValid {
		return NewErrPasswordIsInvalid(sessionId)
	}
	userSession := domain.NewUserSession(user.Id, sessionId)
	c.userSessionRepository.Add(userSession)
	return nil
}

func (c ChatService) ChangePassword(sessionId, oldPassword, newPassword string) error {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionId(sessionId)
	if !userSessionExists {
		return NewErrSessionNotLoggedIn(sessionId)
	}
	user, userExists := c.userRepository.FindById(userSession.UserId)
	if !userExists {
		return errors.New(fmt.Sprintf("user was not found, userId: %s", userSession.UserId))
	}
	if !user.PasswordIsValid(oldPassword) {
		return NewErrPasswordIsInvalid(sessionId)
	}
	err := user.SetPassword(newPassword)
	if err != nil {
		return NewErrPasswordIsInvalid(sessionId)
	}
	return nil
}

func (c ChatService) GetUserNameForSessionId(sessionId string) string {
	userSession, userSessionExists := c.userSessionRepository.FindBySessionId(sessionId)
	if !userSessionExists {
		return ""
	}
	user, userExists := c.userRepository.FindById(userSession.UserId)
	if !userExists {
		return ""
	}
	return user.Name
}

func (c ChatService) GetAllLoggedInUserNames() []string {
	userNames := make([]string, 0)
	for _, user := range c.userRepository.GetAll() {
		userSessions := c.userSessionRepository.DeleteByUserId(user.Id)
		if len(userSessions) > 0 {
			userNames = append(userNames, user.Name)
		}
	}
	return userNames
}

func (c ChatService) QuitSession(sessionId string) {
	session, sessionExists := c.sessionRepository.FindById(sessionId)
	if !sessionExists {
		return
	}
	session.Close <- struct{}{}
	c.userSessionRepository.DeleteBySessionId(sessionId)
	c.sessionRepository.Delete(sessionId)
}
