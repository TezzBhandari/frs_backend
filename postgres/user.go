package postgres

import (
	"context"

	"github.com/TezzBhandari/frs"
)

var _ frs.UserService = (*UserService)(nil)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *frs.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	err = createUser(ctx, tx, user)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// return NOTFOUND Error | UNAUTHORIZED Error
func (s *UserService) FindUserById(ctx context.Context, id int) (*frs.User, error) {
	return nil, nil
}
func (s *UserService) UpdateUser(ctx context.Context, id int, upd frs.UserUpdate) (*frs.User, error) {
	return nil, nil
}

// return NOTFOUND | UNAUTHORIZED Error
func (s *UserService) FindUser(ctx context.Context, filter *frs.FilterUser) ([]*frs.User, int, error) {
	return nil, 0, nil
}

func createUser(ctx context.Context, tx *Tx, user *frs.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}

	user.CreatedAt = tx.Now
	user.ID = tx.db.snowflake.Generate().Int64()

	insertUserQuery := `
		INSERT INTO users (
			id,
			username,
			email,
			password,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err = tx.Exec(ctx, insertUserQuery, user.ID, user.Username, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
