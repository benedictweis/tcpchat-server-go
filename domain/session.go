package domain

import "github.com/google/uuid"

// Session represents a newly created session.
type Session struct {
	ID                string
	MessagesToSession chan<- string
	Close             chan<- interface{}
}

func NewSession(messagesToSession chan<- string, close chan<- interface{}) *Session {
	return &Session{uuid.New().String(), messagesToSession, close}
}

type SessionRepository interface {
	Add(Session) bool
	FindByID(string) (session Session, sessionExists bool)
	FindAllExceptBySessionID(string) []Session
	Delete(string) (session Session, sessionExists bool)
}

type InMemorySessionRepository struct {
	sessions map[string]Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]Session)}
}

func (i *InMemorySessionRepository) Add(session Session) bool {
	if _, sessionExists := i.sessions[session.ID]; sessionExists {
		return false
	}
	i.sessions[session.ID] = session
	return true
}

func (i *InMemorySessionRepository) Delete(sessionID string) (session Session, ok bool) {
	if session, ok = i.sessions[sessionID]; !ok {
		return
	}
	delete(i.sessions, sessionID)
	return
}

func (i *InMemorySessionRepository) FindByID(id string) (session Session, ok bool) {
	session, ok = i.sessions[id]
	return
}

func (i *InMemorySessionRepository) FindAllExceptBySessionID(id string) []Session {
	sessions := make([]Session, 0)
	for _, session := range i.sessions {
		if session.ID != id {
			sessions = append(sessions, session)
		}
	}
	return sessions
}
