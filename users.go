package eauth

import (
	"fmt"
)

// User is the server's user struct. Sessions are attached to users.
// Id is included so urls do not need to include the email in the link url.
// Tokens are refreshed everytime a new session is created in order to prevent
// replay attacks with the given link URL.
type User struct {
	Id    int64
	Email string
	Token string
}

// String returns a string representation of the user id and email
func (u User) String() string {
	return fmt.Sprintf("%d: %s", u.Id, u.Email)
}

// UserManager is used by the server as the interface a backend must implement.
// TODO Save should likely return a User or it's impossible to know
// any auto-create properties

type UserManager interface {
	Save(user User) error
	Delete(user User) error
	UpdateToken(user User, token string) error
	Get(id int64) User
	GetEmail(email string) User
}
