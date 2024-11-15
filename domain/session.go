package domain

import "github.com/google/uuid"

// Session represents a newly created session
type Session struct {
	Id                string
	MessagesToSession chan<- string
	Close             chan<- interface{}
}

func NewSession(messagesToSession chan<- string, close chan<- interface{}) *Session {
	return &Session{uuid.New().String(), messagesToSession, close}
}

type SessionRepository interface {
	Add(Session) bool
	FindById(string) (session Session, sessionExists bool)
	FindAllExceptBySessionId(string) []Session
	Delete(string) (session Session, sessionExists bool)
}

type InMemorySessionRepository struct {
	sessions map[string]Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]Session)}
}

func (i *InMemorySessionRepository) Add(session Session) bool {
	if _, sessionExists := i.sessions[session.Id]; sessionExists {
		return false
	}
	i.sessions[session.Id] = session
	return true
}

func (i *InMemorySessionRepository) Delete(sessionId string) (session Session, ok bool) {
	if session, ok = i.sessions[sessionId]; !ok {
		return
	}
	delete(i.sessions, sessionId)
	return
}

func (i *InMemorySessionRepository) FindById(id string) (session Session, ok bool) {
	session, ok = i.sessions[id]
	return
}

func (i *InMemorySessionRepository) FindAllExceptBySessionId(id string) []Session {
	sessions := make([]Session, 0)
	for _, session := range i.sessions {
		if session.Id != id {
			sessions = append(sessions, session)
		}
	}
	return sessions
}
