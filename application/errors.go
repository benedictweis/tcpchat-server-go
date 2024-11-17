package application

import "fmt"

type UserFriendlyError interface {
	error
	UserFriendlyError() string
}

type BaseError struct {
	sessionId string
	message   string
	userMsg   string
}

func NewBaseError(sessionId, message, userMsg string) BaseError {
	return BaseError{
		sessionId: sessionId,
		message:   message,
		userMsg:   userMsg,
	}
}

func (e BaseError) Error() string {
	return e.message
}

func (e BaseError) UserFriendlyError() string {
	return e.userMsg
}

type ErrSessionNotLoggedIn struct {
	BaseError
}

func NewErrSessionNotLoggedIn(sessionId string) *ErrSessionNotLoggedIn {
	return &ErrSessionNotLoggedIn{NewBaseError(
		sessionId,
		fmt.Sprintf("session %s not logged in", sessionId),
		"you are not logged in",
	)}
}

type ErrMessagePartnerDoesNotExist struct {
	BaseError
}

func NewErrMessagePartnerDoesNotExist(sessionId string, messagePartnerUserName string) *ErrMessagePartnerDoesNotExist {
	return &ErrMessagePartnerDoesNotExist{NewBaseError(
		sessionId,
		fmt.Sprintf("session %s tried to Message non existant partner %s", sessionId, messagePartnerUserName),
		"your Message partner does not seem to be logged in",
	)}
}

type ErrMessagePartnerNotLoggedIn struct {
	BaseError
}

func NewErrMessagePartnerNotLoggedIn(sessionId string, messagePartnerUserName string) *ErrMessagePartnerNotLoggedIn {
	return &ErrMessagePartnerNotLoggedIn{NewBaseError(
		sessionId,
		fmt.Sprintf("session %s tried to Message non logged in partner %s", sessionId, messagePartnerUserName),
		"your Message partner does not seem to be logged in",
	)}
}

type ErrCouldNotCreateUser struct {
	BaseError
}

func NewErrCouldNotCreateUser(sessionId string) *ErrCouldNotCreateUser {
	return &ErrCouldNotCreateUser{NewBaseError(
		sessionId,
		fmt.Sprintf("could not create user for session id: %s", sessionId),
		"could not create user, password is likely invalid",
	)}
}

type ErrUserNameAlreadyExists struct {
	BaseError
}

func NewErrUserNameAlreadyExists(sessionId string, userName string) *ErrUserNameAlreadyExists {
	return &ErrUserNameAlreadyExists{NewBaseError(
		sessionId,
		fmt.Sprintf("session %s tried to create user %s that already exists", sessionId, userName),
		"a user with that name already exists",
	)}
}

type ErrUserDoesNotExist struct {
	BaseError
}

func NewErrUserDoesNotExist(sessionId string, userName string) *ErrUserDoesNotExist {
	return &ErrUserDoesNotExist{NewBaseError(
		sessionId,
		fmt.Sprintf("session %s tried to access user %s that does not exist", sessionId, userName),
		"a user with that name does not exist",
	)}
}

type ErrPasswordIsInvalid struct {
	BaseError
}

func NewErrPasswordIsInvalid(sessionId string) *ErrPasswordIsInvalid {
	return &ErrPasswordIsInvalid{NewBaseError(
		sessionId,
		fmt.Sprintf("session %s entered invalid password", sessionId),
		"wrong password",
	)}
}
