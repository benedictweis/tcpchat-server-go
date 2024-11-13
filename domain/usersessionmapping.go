package domain

type UserSession struct {
	SessionId string
	UserName  string
}

func NewUserSession(sessionId string, userName string) *UserSession {
	return &UserSession{sessionId, userName}
}

type UserSessionRepository interface {
	GetAll() []*UserSession
	Add(*UserSession)
	FindBySessionId(string) (*UserSession, bool)
	FindByUserName(string) []*UserSession
	DeleteBySessionId(string) (*UserSession, bool)
	DeleteByUserName(string) ([]*UserSession, bool)
}

type InMemoryUserSessionRepository struct {
	userSessions map[string]*UserSession
}

func NewInMemoryUserSessionRepository() *InMemoryUserSessionRepository {
	return &InMemoryUserSessionRepository{userSessions: make(map[string]*UserSession)}
}

func (i *InMemoryUserSessionRepository) GetAll() []*UserSession {
	allUserSessions := make([]*UserSession, 0, len(i.userSessions))
	for _, userSession := range i.userSessions {
		allUserSessions = append(allUserSessions, userSession)
	}
	return allUserSessions
}

func (i *InMemoryUserSessionRepository) Add(session *UserSession) {
	i.userSessions[session.SessionId] = session
}

func (i *InMemoryUserSessionRepository) FindBySessionId(sessionId string) (userSession *UserSession, ok bool) {
	userSession, ok = i.userSessions[sessionId]
	return
}

func (i *InMemoryUserSessionRepository) FindByUserName(userName string) []*UserSession {
	userSessions := make([]*UserSession, 0)
	for _, userSession := range i.userSessions {
		if userSession.UserName == userName {
			userSessions = append(userSessions, userSession)
		}
	}
	return userSessions
}

func (i *InMemoryUserSessionRepository) DeleteBySessionId(sessionId string) (userSession *UserSession, ok bool) {
	if userSession, ok = i.userSessions[sessionId]; !ok {
		return
	}
	delete(i.userSessions, sessionId)
	return
}

func (i *InMemoryUserSessionRepository) DeleteByUserName(userName string) []*UserSession {
	userSessions := i.FindByUserName(userName)
	for _, userSession := range i.userSessions {
		_, _ = i.DeleteBySessionId(userSession.SessionId)
	}
	return userSessions
}
