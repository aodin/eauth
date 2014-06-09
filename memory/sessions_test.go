package memory

import (
	"github.com/aodin/eauth"
	"testing"
)

func TestSessions(t *testing.T) {
	// Create an in-memory session manager
	memory := Sessions()
	s := eauth.Session{Key: "yabbadabbado"}

	// Save the session
	if err := memory.Save(s); err != nil {
		t.Fatalf("Could not create session: %s", err)
	}
	if len(memory.keys) != 1 {
		t.Errorf("Unexpected length of sessions: %d", len(memory.keys))
	}
}
