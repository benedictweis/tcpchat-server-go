package handlers_test

import (
	"github.com/benedictweis/tcpchat-server-go/application/handlers"
	"github.com/benedictweis/tcpchat-server-go/domain"
	"github.com/benedictweis/tcpchat-server-go/test"
	mock_application "github.com/benedictweis/tcpchat-server-go/test/mock"
	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/mock/gomock"
)

var _ = Describe("Textmessagehandler", func() {
	Context("#HandleNewSession", func() {
		var (
			ctrl        *gomock.Controller
			textMessage *domain.TextMessage
			chatService *mock_application.MockChatService
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			textMessage = domain.NewTextMessage(test.SESSION_ID_A, test.TEXT_MESSAGE_A)
			chatService = mock_application.NewMockChatService(ctrl)
		})

		Context("when sending a text message", func() {
			It("should send the text message to everyone", func() {
				chatService.EXPECT().SendTextMessageToEveryone(textMessage.SessionID, textMessage.Message).Times(1)
				handlers.HandleTextMessage(*textMessage, chatService)
			})
		})
	})
})
