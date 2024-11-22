package domain_test

import (
	"github.com/benedictweis/tcpchat-server-go/domain"
	"github.com/benedictweis/tcpchat-server-go/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {

	var (
		userNameA           string
		userPasswordA       string
		userPasswordB       string
		userPasswordTooLong string
	)

	BeforeEach(func() {
		userNameA = test.USER_NAME_A
		userPasswordA = test.USER_PASSWORD_A
		userPasswordB = test.USER_PASSWORD_B
		userPasswordTooLong = test.USER_PASSWORD_X_TOO_LONG
	})

	Context("#User", func() {
		Context("#NewUser", func() {
			Context("when presented with a valid password", func() {
				It("should create a valid new user", func() {
					user, err := domain.NewUser(userNameA, userPasswordA)
					Expect(user.ID).To(Not(BeEmpty()))
					Expect(user.Name).To(Equal(userNameA))
					Expect(user.PasswordIsValid(userPasswordA)).To(BeTrue())
					Expect(err).To(BeNil())
				})
			})
			Context("when presented with an invalid password", func() {
				It("should return an error", func() {
					user, err := domain.NewUser(userNameA, userPasswordTooLong)
					Expect(user).To(BeNil())
					Expect(err.Error()).To(ContainSubstring("password length"))
				})
			})
			Context("when changing the password of a user", func() {
				It("should change the password accordingly", func() {
					user, err := domain.NewUser(userNameA, userPasswordA)
					Expect(err).To(BeNil())
					Expect(user.PasswordIsValid(userPasswordA)).To(BeTrue())

					err = user.SetPassword(userPasswordB)
					Expect(err).To(BeNil())
					Expect(user.PasswordIsValid(userPasswordB)).To(BeTrue())
				})
			})
		})
	})
})
