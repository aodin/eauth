package eauth

import (
	"fmt"
)

// TODO Is Id even needed?
type User struct {
	Id    int64
	Email string
	Token string
}

func (u User) String() string {
	return fmt.Sprintf("%d: %s", u.Id, u.Email)
}

type UserManager interface {
	Create(user User) error
	Get(email string) User
	GetId(id int64) User
}
