package domain_test

import (
	"github.com/benedictweis/tcpchat-server-go/domain"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command", func() {
	Context("#CommandType", func() {
		DescribeTable("Getting a Command Type from a string",
			func(commandTypeAsString string, expectedCommandType domain.CommandType) {
				commandType := domain.CommandTypeFromString(commandTypeAsString)
				Expect(commandType).To(Equal(expectedCommandType))
			},
			// Valid commands
			Entry("When given valid command unknown", "unknown", domain.Unknown),
			Entry("When given valid command name", "name", domain.ChangeName),
			Entry("When given valid command msg", "msg", domain.PrivateMessage),
			Entry("When given valid command acc", "acc", domain.CreateAccount),
			Entry("When given valid command login", "login", domain.Login),
			Entry("When given valid command passwd", "passwd", domain.ChangePassword),
			Entry("When given valid command info", "info", domain.Info),
			Entry("When given valid command who", "who", domain.Who),
			Entry("When given valid command quit", "quit", domain.Quit),
			// Invalid Commands
			Entry("When given invalid command <empty string>", "", domain.Unknown),
			Entry("When given invalid command 1234", "1234", domain.Unknown),
			Entry("When given invalid command msg1", "msg1", domain.Unknown),
			Entry("When given invalid command does not exist", "does not exist", domain.Unknown),
			Entry("When given invalid command account", "account", domain.Unknown),
		)
		DescribeTable("Getting a String from a CommandType",
			func(commandType domain.CommandType, expectedString string) {
				commandTypeAsString := commandType.String()
				Expect(commandTypeAsString).To(Equal(expectedString))
			},
			// Valid command types
			Entry("When given valid CommandType Unknown", domain.Unknown, "unknown"),
			Entry("When given valid CommandType ChangeName", domain.ChangeName, "name"),
			Entry("When given valid CommandType PrivateMessage", domain.PrivateMessage, "msg"),
			Entry("When given valid CommandType CreateAccount", domain.CreateAccount, "acc"),
			Entry("When given valid CommandType Login", domain.Login, "login"),
			Entry("When given valid CommandType ChangePassword", domain.ChangePassword, "passwd"),
			Entry("When given valid CommandType Info", domain.Info, "info"),
			Entry("When given valid CommandType Who", domain.Who, "who"),
			Entry("When given valid CommandType Quit", domain.Quit, "quit"),
			// Invalid command types
			Entry("When given invalid CommandType 9", domain.CommandType(9), "9"),
			Entry("When given invalid CommandType 500", domain.CommandType(500), "500"),
			Entry("When given invalid CommandType -1", domain.CommandType(-1), "-1"),
			Entry("When given invalid CommandType 1000000", domain.CommandType(1000000), "1000000"),
		)
	})
})
