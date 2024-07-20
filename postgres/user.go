package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/TezzBhandari/frs"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
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
	// don't use pointer for txOptions to make it optional, does not return a transation
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
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
func (s *UserService) FindUsers(ctx context.Context, filterUser *frs.FilterUser) ([]*frs.User, int, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	users, n, err := findUsers(ctx, tx, filterUser)
	if err != nil {
		return nil, 0, err
	}
	return users, n, nil
}

func createUser(ctx context.Context, tx *Tx, user *frs.User) error {
	log.Debug().Msg("reached created user")
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

func findUsers(ctx context.Context, tx *Tx, filterUser *frs.FilterUser) ([]*frs.User, int, error) {
	where := []string{"1 = 1"}
	args := []any{}
	var i int = 1

	if filterUser.Username != nil {
		where = append(where, fmt.Sprintf("username = $%d", i))
		args = append(args, *filterUser.Username)
		i++
	}

	if filterUser.Email != nil {
		where = append(where, fmt.Sprintf("email = $%d", i))
		args = append(args, *filterUser.Email)
		i++
	}

	whereClause := strings.Join(where, " AND ")
	findUserQuery := `SELECT id, username, email, created_at FROM users WHERE ` + whereClause

	log.Debug().Msg(findUserQuery)

	rows, err := tx.Query(ctx, findUserQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var users []*frs.User
	var user frs.User

	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	return users, len(users), nil
}
