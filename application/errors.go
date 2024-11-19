// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package application

import "fmt"

type UserFriendlyError interface {
	error
	UserFriendlyError() string
}

type BaseError struct {
	sessionID string
	message   string
	userMsg   string
}

func NewBaseError(sessionID, message, userMsg string) BaseError {
	return BaseError{
		sessionID: sessionID,
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

func NewErrSessionNotLoggedIn(sessionID string) *ErrSessionNotLoggedIn {
	return &ErrSessionNotLoggedIn{NewBaseError(
		sessionID,
		fmt.Sprintf("session %s not logged in", sessionID),
		"you are not logged in",
	)}
}

type ErrMessagePartnerDoesNotExist struct {
	BaseError
}

func NewErrMessagePartnerDoesNotExist(sessionID string, messagePartnerUserName string) *ErrMessagePartnerDoesNotExist {
	return &ErrMessagePartnerDoesNotExist{NewBaseError(
		sessionID,
		fmt.Sprintf("session %s tried to Message non existant partner %s", sessionID, messagePartnerUserName),
		"your Message partner does not seem to be logged in",
	)}
}

type ErrMessagePartnerNotLoggedIn struct {
	BaseError
}

func NewErrMessagePartnerNotLoggedIn(sessionID string, messagePartnerUserName string) *ErrMessagePartnerNotLoggedIn {
	return &ErrMessagePartnerNotLoggedIn{NewBaseError(
		sessionID,
		fmt.Sprintf("session %s tried to Message non logged in partner %s", sessionID, messagePartnerUserName),
		"your Message partner does not seem to be logged in",
	)}
}

type ErrCouldNotCreateUser struct {
	BaseError
}

func NewErrCouldNotCreateUser(sessionID string) *ErrCouldNotCreateUser {
	return &ErrCouldNotCreateUser{NewBaseError(
		sessionID,
		fmt.Sprintf("could not create user for session id: %s", sessionID),
		"could not create user, password is likely invalid",
	)}
}

type ErrUserNameAlreadyExists struct {
	BaseError
}

func NewErrUserNameAlreadyExists(sessionID string, userName string) *ErrUserNameAlreadyExists {
	return &ErrUserNameAlreadyExists{NewBaseError(
		sessionID,
		fmt.Sprintf("session %s tried to create user %s that already exists", sessionID, userName),
		"a user with that name already exists",
	)}
}

type ErrUserDoesNotExist struct {
	BaseError
}

func NewErrUserDoesNotExist(sessionID string, userName string) *ErrUserDoesNotExist {
	return &ErrUserDoesNotExist{NewBaseError(
		sessionID,
		fmt.Sprintf("session %s tried to access user %s that does not exist", sessionID, userName),
		"a user with that name does not exist",
	)}
}

type ErrPasswordIsInvalid struct {
	BaseError
}

func NewErrPasswordIsInvalid(sessionID string) *ErrPasswordIsInvalid {
	return &ErrPasswordIsInvalid{NewBaseError(
		sessionID,
		fmt.Sprintf("session %s entered invalid password", sessionID),
		"wrong password",
	)}
}
