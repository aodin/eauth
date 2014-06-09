package eauth

import (
	"fmt"
	"testing"
)

func TestRandomKey(t *testing.T) {
	// Generate a random key and check for errors
	key, err := RandomKey()
	if err != nil {
		t.Fatalf("Error during RandomKey(): %s", err)
	}
	if key == "" {
		t.Fatal("Blank key returned from RandomKey()")
	}
	if len(key) != 24 {
		t.Fatal("Unexpected key length from RandomKey(): %d", len(key))
	}
}

// A key generator that returns the same key once, then a random one
// Used for testing collisions and is not exported
type badKeyGen int

// This must be a pointer method or the count will be passed by value!
func (kg *badKeyGen) KeyFunc() (string, error) {
	if *kg < 1 {
		*kg += 1
		return "BAD", nil
	}

	return RandomKey()
}

// func TestNewSession(t *testing.T) {
// 	// Create an in-memory session managerd
// 	manager := InMemorySessions()

// 	session, err := NewSession(manager, 1, defaultCookie)
// 	if err != nil {
// 		t.Fatalf("Error during NewSession(): %s", err)
// 	}
// 	if session.Key == "" {
// 		t.Error("Blank session key returned from NewSession()")
// 	}
// 	if len(manager) != 1 {
// 		t.Error("Session was not created by NewSession()")
// 	}

// 	// Reset the manager and test with sessions with the bad key generator
// 	manager = InMemorySessions()
// 	var bad badKeyGen
// 	session, err = newSession(manager, bad.KeyFunc, 1, defaultCookie)
// 	if err != nil {
// 		t.Fatalf("Error during newSession(): %s", err)
// 	}

// 	// A session should have been created with the key bad
// 	if _, ok := manager["BAD"]; !ok {
// 		t.Fatal("A bad session key was not created by newSession()")
// 	}

// 	// Reset the counter so that another BAD key is created
// 	bad = 0
// 	session, err = newSession(manager, bad.KeyFunc, 1, defaultCookie)
// 	if err != nil {
// 		t.Fatalf("Error during repeated newSession(): %s", err)
// 	}

// 	// Two sessions should exist
// 	if len(manager) != 2 {
// 		t.Fatalf("Unexpected number of sessions from newSession(): %d", len(manager))
// 	}
// }
