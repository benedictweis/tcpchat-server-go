// Copyright (c) 2024 Benedict Weis. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package domain

// TextMessage represents a message a user intends to send.
type TextMessage struct {
	SessionID string
	Message   string
}

func NewTextMessage(sessionID string, message string) *TextMessage {
	return &TextMessage{SessionID: sessionID, Message: message}
}
