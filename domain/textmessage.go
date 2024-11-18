package domain

// TextMessage represents a message a user intends to send.
type TextMessage struct {
	SessionID string
	Message   string
}
