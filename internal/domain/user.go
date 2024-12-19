package domain

import (
	"time"

	"github.com/lithammer/shortuuid/v4"
	"http/internal/tberrors"
)

type User struct {
	ID        string
	Name      string
	DeletedAt *time.Time
}

func NewUser(name string) (*User, error) {
	u := &User{
		ID:   shortuuid.New(),
		Name: name,
	}

	return u.validate()
}

var invalidEmptyNameError = tberrors.NewValidationError("invalid empty name", "name")

func (u *User) validate() (*User, error) {
	if u.Name == "" {
		return nil, invalidEmptyNameError
	}

	return u, nil
}
