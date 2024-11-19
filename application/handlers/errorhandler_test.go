// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package handlers

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
)

// TODO add test for error handler

var _ = Describe("Error Handler", func() {
	Describe("#handleErrors", func() {
		Context("when the error is a UserFriendlyError", func() {
			It("should be sent to the session as a UserFriendlyError", func() {
				handleErrors(fmt.Errorf("hi, this is an error"), nil, "1234")
			})
		})

		Context("when the error is a not UserFriendlyError", func() {
			It("should be a sent to the session as an internal server error", func() {
			})
		})
	})
})
