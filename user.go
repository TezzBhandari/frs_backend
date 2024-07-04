package frs

import "context"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type FilterUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService interface {
	createUser(ctx context.Context, user *User) error
	// return NOTFOUND Error | UNAUTHORIZED Error
	findUserById(ctx context.Context, id int) (*User, error)
	updateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)
	// return NOTFOUND | UNAUTHORIZED Error
	findUser(ctx context.Context, filter *FilterUser) ([]*User, int, error)
}
