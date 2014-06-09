package memory

import (
	"github.com/aodin/eauth"
	"testing"
)

func TestUsers(t *testing.T) {
	memory := Users()
	admin := eauth.User{Email: "admin@example.com"}

	// Create a User
	if err := memory.Save(admin); err != nil {
		t.Fatalf("Could not create user: %s", err)
	}
	if len(memory.ids) != 1 {
		t.Errorf("Unexpected length of user ids: %d", len(memory.ids))
	}
	if len(memory.emails) != 1 {
		t.Errorf("Unexpected length of user emails: %d", len(memory.emails))
	}

	// Get a user by email
	u := memory.GetEmail(admin.Email)
	if u.Email != admin.Email {
		t.Errorf("Unexpected user email: %s", u.Email)
	}

	// TODO Delete the user
}
