package domain

// Session represents a newly created session
type Session struct {
	SessionId         string
	MessagesToSession chan<- string
	Close             chan<- interface{}
}

type SessionRepository interface {
	Add(Session)
	Delete(id string) (Session, bool)
	FindById(id string) (Session, bool)
}

type InMemorySessionRepository struct {
	sessions map[string]Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]Session)}
}

func (i *InMemorySessionRepository) Add(session Session) {
	i.sessions[session.SessionId] = session
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
