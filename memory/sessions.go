package memory

import (
	"fmt"
	"github.com/aodin/eauth"
)

// TODO mutexes!
type sessions struct {
	keys map[string]eauth.Session
}

// Save adds a new session to the in-memory store of sessions
func (s *sessions) Save(session eauth.Session) error {
	if _, exists := s.keys[session.Key]; exists {
		return fmt.Errorf(
			"A session with the key %s already exists",
			session.Key,
		)
	}
	s.keys[session.Key] = session
	return nil
}

// Delete will remove the session from the in-memory store of sessions.
// It will error if the given session does not exist.
func (s *sessions) Delete(key string) error {
	if _, exists := s.keys[key]; !exists {
		return fmt.Errorf(
			"No session with the key %s exists",
			key,
		)
	}
	delete(s.keys, key)
	return nil
}

// Get returns the session with the given key. It will return a zero-
// initialized session if the key does not exist.
func (s *sessions) Get(key string) eauth.Session {
	return s.keys[key]
}

// Sessions creates an in-memory store of sessions that implements the methods
// of the eauth.SessionManager interface.
func Sessions() *sessions {
	return &sessions{
		keys: make(map[string]eauth.Session),
	}
}
