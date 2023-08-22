// Package: telnetter
// Description: Defines and manages Telnet sessions and their lifecycle.
// Git Repository: [URL not provided in the source]
// License: [License type not provided in the source]
package telnetter

import (
	"crypto/rand"
	"encoding/hex"
)

// Session represents an individual Telnet session with a unique ID, connection details, and status.
type Session struct {
    ID     string
    Conn   *Conn
    Active bool
}

// SessionManager is responsible for managing and tracking all active sessions.
type SessionManager struct {
    Sessions map[string]*Session
}


// Title: New Session Manager
// Description: Initializes and returns a new SessionManager for managing Telnet sessions.
// Function: func NewSessionManager() *SessionManager
// CalledWith: manager := NewSessionManager()
// ExpectedOutput: A pointer to a new SessionManager instance.
// Example: sessionManager := NewSessionManager()
func NewSessionManager() *SessionManager {
    return &SessionManager{
        Sessions: make(map[string]*Session),
    }
}


// Title: Create New Session
// Description: Generates a unique session ID, creates a new session with the provided connection, and adds it to the session manager.
// Function: func (sm *SessionManager) CreateSession(conn *Conn) *Session
// CalledWith: session := manager.CreateSession(conn)
// ExpectedOutput: A pointer to the created Session instance.
// Example: newSession := sessionManager.CreateSession(conn)
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


// Title: Get Session
// Description: Retrieves a session based on its ID.
// Function: func (sm *SessionManager) GetSession(id string) (*Session, bool)
// CalledWith: session, found := manager.GetSession("sessionId")
// ExpectedOutput: The session and a boolean indicating if the session was found.
// Example: sess, exists := sessionManager.GetSession("12345")
func (sm *SessionManager) GetSession(id string) (*Session, bool) {
    session, found := sm.Sessions[id]
    return session, found
}


// Title: End Session
// Description: Marks a session as inactive based on its ID.
// Function: func (sm *SessionManager) EndSession(id string)
// CalledWith: manager.EndSession("sessionId")
// ExpectedOutput: None, marks the session as inactive.
// Example: sessionManager.EndSession("12345")
func (sm *SessionManager) EndSession(id string) {
    if session, found := sm.Sessions[id]; found {
        session.Active = false
    }
}

// generateSessionID creates a unique session ID using random bytes and returns it as a hexadecimal string.
func generateSessionID() string {
    b := make([]byte, 16)
    rand.Read(b)
    return hex.EncodeToString(b)
}