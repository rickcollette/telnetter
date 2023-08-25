package telnetter

import (
	"crypto/rand"
	"encoding/hex"
)

type Session struct {
    ID     string
    Conn   *Conn
    Active bool
}

type SessionManager struct {
    Sessions map[string]*Session
}



func NewSessionManager() *SessionManager {
    return &SessionManager{
        Sessions: make(map[string]*Session),
    }
}



func (sm *SessionManager) CreateSession(conn *Conn) *Session {
    id := generateSessionID()
    session := &Session{
        ID:     id,
        Conn:   conn,
        Active: true,
    }
    sm.Sessions[id] = session
    return session
}



func (sm *SessionManager) GetSession(id string) (*Session, bool) {
    session, found := sm.Sessions[id]
    return session, found
}


func (sm *SessionManager) EndSession(id string) {
    if session, found := sm.Sessions[id]; found {
        session.Active = false
    }
}

func generateSessionID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}