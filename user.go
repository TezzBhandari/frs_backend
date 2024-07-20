package frs

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) FromJson(v io.ReadCloser) error {
	return json.NewDecoder(v).Decode(u)
}

func (u *User) Validate() error {
	if u.Username == "" {
		return Errorf(EBADREQUEST, "username required")
	}

	if u.Email == "" {
		return Errorf(EBADREQUEST, "email required")
	}

	if validEmail(u.Email) {
		return Errorf(EBADREQUEST, "invalid email")
	}

	if u.Password == "" {
		return Errorf(EBADREQUEST, "password required")
	}

	if len(u.Password) < 8 {
		return Errorf(EBADREQUEST, "password should be at least 8 character long")
	}

	return nil
}

type UserUpdate struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type FilterUser struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	// return NOTFOUND Error | UNAUTHORIZED Error
	FindUserById(ctx context.Context, id int) (*User, error)
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)
	// return NOTFOUND | UNAUTHORIZED Error
	FindUsers(ctx context.Context, filter *FilterUser) ([]*User, int, error)
}

func validEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// func validPassword(password string) bool {
// 	passwordRegex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$`
// 	re := regexp.MustCompile(passwordRegex)
// 	return re.MatchString(password)
// }
