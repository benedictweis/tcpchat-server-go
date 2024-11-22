// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package handlers

import (
	"fmt"
	"github.com/benedictweis/tcpchat-server-go/application"
	"github.com/benedictweis/tcpchat-server-go/test"
	mock_application "github.com/benedictweis/tcpchat-server-go/test/mock"
	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/mock/gomock"
)

// TODO add test for error handler

var _ = Describe("Error Handler", func() {
	Context("#handleErrors", func() {
		var (
			ctrl                 *gomock.Controller
			sessionID            string
			userFriendlyError    application.UserFriendlyError
			nonUserFriendlyError error
			chatService          *mock_application.MockChatService
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			sessionID = test.SESSION_ID_A
			userFriendlyError = application.NewErrSessionNotLoggedIn(sessionID)
			nonUserFriendlyError = fmt.Errorf(test.UNKNOWN_ERROR_TEXT)
			chatService = mock_application.NewMockChatService(ctrl)
		})

		Context("when the error is a UserFriendlyError", func() {
			It("should be sent to the session as a UserFriendlyError", func() {
				chatService.EXPECT().SendMessageToSessionFromServer(sessionID, userFriendlyError.UserFriendlyError()).Times(1)
				handleErrors(userFriendlyError, chatService, sessionID)
			})
		})

		Context("when the error is a not UserFriendlyError", func() {
			It("should be a sent to the session as an internal server error", func() {
				chatService.EXPECT().SendMessageToSessionFromServer(sessionID, "internal server error").Times(1)
				handleErrors(nonUserFriendlyError, chatService, sessionID)
			})
		})
	})
})
