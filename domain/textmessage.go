package domain

// TextMessage represents a message a user intends to send
type TextMessage struct {
	SessionId string
	Message   string
}
