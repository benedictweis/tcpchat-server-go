package handlers_test

import (
	"github.com/benedictweis/tcpchat-server-go/application/handlers"
	"github.com/benedictweis/tcpchat-server-go/domain"
	mock_application "github.com/benedictweis/tcpchat-server-go/test/mock"
	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/mock/gomock"
)

var _ = Describe("NewSessionHandler", func() {
	Context("#HandleNewSession", func() {
		var (
			ctrl        *gomock.Controller
			session     *domain.Session
			chatService *mock_application.MockChatService
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			session = domain.NewSession(make(chan<- string), make(chan<- interface{}))
			chatService = mock_application.NewMockChatService(ctrl)
		})

		Context("when registering a new session", func() {
			It("should register the session and inform the user", func() {
				chatService.EXPECT().RegisterNewSession(*session).Times(1)
				chatService.EXPECT().SendMessageToSessionFromServer(session.ID, "Welcome to this server!").Times(1)
				handlers.HandleNewSession(*session, chatService)
			})
		})
	})
})
