package eauth

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// Includes options for making sessions even more secure:
// * Single sessions per user
// * IP address fixation
// Session does not include data.
type Session struct {
	Key     string
	UserId  int64
	IP      string
	Expires time.Time
}

type KeyFunc func() (string, error)

// For 144 bit sessions, we'll need to generate 18 random bytes.
// These will be encoded in URL safe base 64, for a length of 24 chars.
func RandomKey() (string, error) {
	b := make([]byte, 18)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Session keys become the cookie's value. US-ASCII is safe except for
// control characters, commas, semicolons and backslash.
// URL-encoded base64 is safe and is used here.
func NewSession(m SessionManager, uid int64, c CookieConfig) (Session, error) {
	return newSession(m, RandomKey, uid, c)
}

func newSession(m SessionManager, key KeyFunc, uid int64, c CookieConfig) (Session, error) {
	// Start a new session
	session := Session{
		UserId: uid,
	}
	// TODO Set the expires from the cookie config (if not zero)
	session.Expires = time.Now().Add(c.Age)

	// Generate a new session key
	var err error
	for {
		session.Key, err = key()
		if err != nil {
			return session, err
		}
		// Attempt to get the session with this key
		// A zero-init session will be returned if the key does not exist
		if s := m.Get(session.Key); s.Key == "" {
			break
		}
	}

	if err = m.Save(session); err != nil {
		return session, err
	}
	return session, nil
}

// IsValidSession checks if a session key exists in the given manager.
func IsValidSession(m SessionManager, key string) bool {
	return isValidSession(m, key, time.Now())
}

func isValidSession(m SessionManager, key string, now time.Time) bool {
	// Exploit the fact that Get will return a zero-initialized session if
	// the key is not valid - that will include a zero timestamp Expires.
	return m.Get(key).Expires.After(now)
}

// SessionManager is the persistance interface for sessions.
type SessionManager interface {
	Save(session Session) error
	Delete(key string) error
	Get(key string) Session
}
