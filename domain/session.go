package domain

// Session represents a newly created session
type Session struct {
	Id                string
	MessagesToSession chan<- string
	Close             chan<- interface{}
}

type SessionRepository interface {
	Add(Session)
	Delete(id string) (Session, bool)
	FindById(id string) (Session, bool)
	FindAllExceptBySessionId(id string) []Session
}

type InMemorySessionRepository struct {
	sessions map[string]Session
}

func NewInMemorySessionRepository() *InMemorySessionRepository {
	return &InMemorySessionRepository{sessions: make(map[string]Session)}
}

func (i *InMemorySessionRepository) Add(session Session) {
	i.sessions[session.Id] = session
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
