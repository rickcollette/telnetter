package telnetter

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type Session struct {
    ID     string
    Conn   *Conn
    Active bool
}

type SessionManager struct {
    Sessions map[string]*Session
    mu       sync.Mutex
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
    sm.mu.Lock()
    sm.Sessions[id] = session
    sm.mu.Unlock()
    return session
}

func (sm *SessionManager) GetSession(id string) (*Session, bool) {
    sm.mu.Lock()
    session, found := sm.Sessions[id]
    sm.mu.Unlock()
    return session, found
}

func (sm *SessionManager) EndSession(id string) {
    sm.mu.Lock()
    if session, found := sm.Sessions[id]; found {
        session.Active = false
        delete(sm.Sessions, id)
    }
    sm.mu.Unlock()
}

func (sm *SessionManager) AddSession(conn *Conn) string {
    session := sm.CreateSession(conn)
    return session.ID
}

func (sm *SessionManager) RemoveSession(id string) {
    sm.EndSession(id)
}

func generateSessionID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}
