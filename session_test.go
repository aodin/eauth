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
