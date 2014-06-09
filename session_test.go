package eauth

import (
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
