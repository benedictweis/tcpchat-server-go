package domain

type UserSession struct {
	UserId    string
	SessionId string
}

func NewUserSession(userId string, sessionId string) *UserSession {
	return &UserSession{userId, sessionId}
}

type UserSessionRepository interface {
	Add(*UserSession)
	GetAll() []*UserSession
	FindBySessionId(string) (*UserSession, bool)
	FindByUserId(string) []*UserSession
	DeleteBySessionId(string) (*UserSession, bool)
	DeleteByUserId(string) []*UserSession
}

type InMemoryUserSessionRepository struct {
	userSessions map[string]*UserSession
}

func NewInMemoryUserSessionRepository() *InMemoryUserSessionRepository {
	return &InMemoryUserSessionRepository{userSessions: make(map[string]*UserSession)}
}

func (i *InMemoryUserSessionRepository) Add(session *UserSession) {
	i.userSessions[session.SessionId] = session
}

func (i *InMemoryUserSessionRepository) GetAll() []*UserSession {
	allUserSessions := make([]*UserSession, 0, len(i.userSessions))
	for _, userSession := range i.userSessions {
		allUserSessions = append(allUserSessions, userSession)
	}
	return allUserSessions
}

func (i *InMemoryUserSessionRepository) FindBySessionId(sessionId string) (userSession *UserSession, ok bool) {
	userSession, ok = i.userSessions[sessionId]
	return
}

func (i *InMemoryUserSessionRepository) FindByUserId(userId string) []*UserSession {
	userSessions := make([]*UserSession, 0)
	for _, userSession := range i.userSessions {
		if userSession.UserId == userId {
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

func (i *InMemoryUserSessionRepository) DeleteByUserId(userId string) []*UserSession {
	userSessions := i.FindByUserId(userId)
	for _, userSession := range i.userSessions {
		_, _ = i.DeleteBySessionId(userSession.SessionId)
	}
	return userSessions
}
