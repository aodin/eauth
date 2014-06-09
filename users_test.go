package eauth

import (
	"fmt"
	"testing"
)

// TODO How to get by email or user id?
// TODO mutexes!
type userMap struct {
	users  []User
	emails map[string]User
}

// Return a zero-initialized user if it does not exist
func (u *userMap) Get(email string) User {
	return u.emails[email]
}

func (u *userMap) GetId(id int64) User {
	intId := int(id)
	if intId >= len(u.users) {
		return User{}
	}
	return u.users[intId]
}

func (u *userMap) Create(user User) error {
	if _, exists := u.emails[user.Email]; exists {
		return fmt.Errorf(
			"A user with the email %s already exists",
			user.Email,
		)
	}
	// Set the user id
	// TODO This means there's no way to delete users!
	user.Id = int64(len(u.users) + 1)
	u.users = append(u.users, user)
	u.emails[user.Email] = user
	fmt.Println("Users:", u.users, u.emails)
	return nil
}

// TODO Bootstrap with users
func InMemoryUsers() *userMap {
	return &userMap{
		users:  make([]User, 0),
		emails: make(map[string]User),
	}
}

func TestUserMap(t *testing.T) {
	memory := InMemoryUsers()
	admin := User{Id: 1, Email: "admin@example.com"}
	if err := memory.Create(admin); err != nil {
		t.Fatalf("Could not create user in memory store: %s", err)
	}
	if len(memory.users) != 1 {
		t.Errorf("Unexpected length of users in memory: %d", len(memory.users))
	}
}
