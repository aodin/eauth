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
	UserId  int64
	Key     string
	IP      string
	Expires time.Time
}

// For 144 bit sessions, we'll need to generate 18 random bytes.
// These will be encoded in URL safe base 64, for a length of 24 chars.
// Also with all io
func RandomKey() (string, error) {
	b := make([]byte, 18)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func KeyExists() bool {
	return false
}

// Session keys become the cookie's value. US-ASCII is safe except for
// control characters, commas, semicolons and backslash.
// URL-encoded base64 is safe and is used here.
func NewSession(userId int64) (Session, error) {
	return Session{}, nil
}
