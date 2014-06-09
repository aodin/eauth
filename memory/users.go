package memory

import (
	"fmt"
	"github.com/aodin/eauth"
)

// TODO mutexes!
type users struct {
	c      int64
	ids    map[int64]eauth.User
	emails map[string]eauth.User
}

// Save adds a new user to the in-memory store of users. It will
// set the Id, overwriting any previous id. Users must have unique emails.
func (u *users) Save(user eauth.User) error {
	// TODO Confirm the email is valid
	if user.Email == "" {
		return fmt.Errorf("Emails are required")
	}
	if _, exists := u.emails[user.Email]; exists {
		return fmt.Errorf(
			"A user with the email %s already exists",
			user.Email,
		)
	}
	// Set the user id
	u.c += 1
	user.Id = u.c

	// Add the user to the id and email maps
	u.ids[user.Id] = user
	u.emails[user.Email] = user
	return nil
}

// Delete will remove the user from the in-memory store of users. It will
// error if the given user does not exist.
func (u *users) Delete(user eauth.User) error {
	if _, exists := u.ids[user.Id]; !exists {
		return fmt.Errorf(
			"No user with the id %d exists",
			user.Id,
		)
	}
	if _, exists := u.emails[user.Email]; !exists {
		return fmt.Errorf(
			"No user with the email %s exists",
			user.Email,
		)
	}
	delete(u.ids, user.Id)
	delete(u.emails, user.Email)
	return nil
}

// Get returns the user with the given id. It will return a zero-initialized
// user if the id does not exist.
func (u *users) Get(id int64) eauth.User {
	return u.ids[id]
}

// GetEmail returns the user with the given email. It will return a zero-
// initialized user if the email does not exist.
func (u *users) GetEmail(email string) eauth.User {
	return u.emails[email]
}

func (u *users) UpdateToken(user eauth.User, token string) error {
	uid, exists := u.ids[user.Id]
	if !exists {
		return fmt.Errorf(
			"No user with the id %d exists",
			user.Id,
		)
	}
	ue, exists := u.emails[user.Email]
	if !exists {
		return fmt.Errorf(
			"No user with the email %s exists",
			user.Email,
		)
	}

	// The users should match
	if uid.Email != user.Email || ue.Id != user.Id {
		return fmt.Errorf(
			"Data integrity error in stored users for user %s",
			user,
		)
	}

	// Update each map with the new token
	user.Token = token
	u.ids[user.Id] = user
	u.emails[user.Email] = user
	return nil
}

// Users creates an in-memory store of users that implements the methods
// of the eauth.UserManager interface.
// TODO Bootstrap with users
func Users() *users {
	return &users{
		ids:    make(map[int64]eauth.User),
		emails: make(map[string]eauth.User),
	}
}
