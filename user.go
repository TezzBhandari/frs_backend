package frs

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	if u.Username == "" {
		return Errorf(EBADREQUEST, "username required")
	}

	if u.Email == "" {
		return Errorf(EBADREQUEST, "email required")
	}

	if u.Password == "" {
		return Errorf(EBADREQUEST, "password required")
	}

	return nil
}

type UserUpdate struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type FilterUser struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	// return NOTFOUND Error | UNAUTHORIZED Error
	FindUserById(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)
	// return NOTFOUND | UNAUTHORIZED Error
	FindUser(ctx context.Context, filter *FilterUser) ([]*User, int, error)
}
