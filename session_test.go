package telnetter

import (
	"testing"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()
	if sm == nil {
		t.Fatal("failed to create a new session manager")
	}

	if len(sm.Sessions) != 0 {
		t.Errorf("expected 0 sessions, got %d", len(sm.Sessions))
	}
}

func TestSessionManager_CreateSession(t *testing.T) {
	sm := NewSessionManager()
	conn := &Conn{}
	session := sm.CreateSession(conn)

	if session == nil {
		t.Fatal("failed to create a new session")
	}

	if session.Conn != conn {
		t.Error("connection in session doesn't match the original connection")
	}

	if _, exists := sm.Sessions[session.ID]; !exists {
		t.Error("session not found in session manager's sessions")
	}
}

func TestSessionManager_GetSession(t *testing.T) {
	sm := NewSessionManager()
	conn := &Conn{}
	session := sm.CreateSession(conn)

	retSession, exists := sm.GetSession(session.ID)
	if !exists {
		t.Error("session not found in session manager's sessions")
	}

	if retSession.ID != session.ID {
		t.Errorf("expected session ID %s, got %s", session.ID, retSession.ID)
	}
}

func TestSessionManager_EndSession(t *testing.T) {
	sm := NewSessionManager()
	conn := &Conn{}
	session := sm.CreateSession(conn)

	sm.EndSession(session.ID)
	if session.Active {
		t.Error("expected session to be inactive, but it's still active")
	}
}
