package domain

type UserSession struct {
	UserID    string
	SessionID string
}

func NewUserSession(userID string, sessionID string) *UserSession {
	return &UserSession{userID, sessionID}
}

type UserSessionRepository interface {
	Add(*UserSession)
	GetAll() []*UserSession
	FindBySessionID(string) (*UserSession, bool)
	FindByUserID(string) []*UserSession
	DeleteBySessionID(string) (*UserSession, bool)
	DeleteByUserID(string) []*UserSession
}

type InMemoryUserSessionRepository struct {
	userSessions map[string]*UserSession
}

func NewInMemoryUserSessionRepository() *InMemoryUserSessionRepository {
	return &InMemoryUserSessionRepository{userSessions: make(map[string]*UserSession)}
}

func (i *InMemoryUserSessionRepository) Add(session *UserSession) {
	i.userSessions[session.SessionID] = session
}

func (i *InMemoryUserSessionRepository) GetAll() []*UserSession {
	allUserSessions := make([]*UserSession, 0, len(i.userSessions))
	for _, userSession := range i.userSessions {
		allUserSessions = append(allUserSessions, userSession)
	}
	return allUserSessions
}

func (i *InMemoryUserSessionRepository) FindBySessionID(sessionID string) (userSession *UserSession, ok bool) {
	userSession, ok = i.userSessions[sessionID]
	return
}

func (i *InMemoryUserSessionRepository) FindByUserID(userID string) []*UserSession {
	userSessions := make([]*UserSession, 0)
	for _, userSession := range i.userSessions {
		if userSession.UserID == userID {
			userSessions = append(userSessions, userSession)
		}
	}
	return userSessions
}

func (i *InMemoryUserSessionRepository) DeleteBySessionID(sessionID string) (userSession *UserSession, ok bool) {
	if userSession, ok = i.userSessions[sessionID]; !ok {
		return
	}
	delete(i.userSessions, sessionID)
	return
}

func (i *InMemoryUserSessionRepository) DeleteByUserID(userID string) []*UserSession {
	userSessions := i.FindByUserID(userID)
	for _, userSession := range i.userSessions {
		_, _ = i.DeleteBySessionID(userSession.SessionID)
	}
	return userSessions
}
