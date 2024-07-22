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

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = findUserById(ctx, tx, id)

	if err != nil {
		return err
	}

	err = deleteUser(ctx, tx, id)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
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
func (s *UserService) FindUserById(ctx context.Context, id int64) (*frs.User, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	user, err := findUserById(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, updUser frs.UpdateUser) (*frs.User, error) {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	user, err := findUserById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if *updUser.Email == "" {
		updUser.Email = &user.Email

	}

	if *updUser.Username == "" {
		updUser.Username = &user.Username

	}

	// TODO: user can only edit only it info. right now any one can edit based on user id

	user, err = updateUser(ctx, tx, id, updUser)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
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
	err := user.Validate()
	if err != nil {
		return err
	}

	user.CreatedAt = tx.Now
	user.ID = int64(tx.db.snowflake.Generate().Int64())
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

	if filterUser.Id != nil {
		where = append(where, fmt.Sprintf("id = $%d", i))
		args = append(args, filterUser.Id)
		i++
	}

	whereClause := strings.Join(where, " AND ")
	findUserQuery := `
	SELECT 
	id, username, email, created_at
	FROM users WHERE ` + whereClause

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

func findUserById(ctx context.Context, tx *Tx, id int64) (*frs.User, error) {
	user, n, err := findUsers(ctx, tx, &frs.FilterUser{Id: &id})
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, frs.Errorf(frs.ENOTFOUND, "user does not exist")
	}
	return user[0], nil
}

func deleteUser(ctx context.Context, tx *Tx, id int64) error {
	deleteUserQuery := `DELETE FROM users WHERE  id = $1`
	_, err := tx.Exec(ctx, deleteUserQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func updateUser(ctx context.Context, tx *Tx, id int64, updateUser frs.UpdateUser) (*frs.User, error) {
	updateUserQuery := `
	UPDATE users
	SET username = $1, email = $2, updated_at = $3
	WHERE id = $4;
	`
	pgTag, err := tx.Exec(ctx, updateUserQuery, *updateUser.Username, *updateUser.Email, tx.Now, id)
	if err != nil {
		return nil, err
	}

	fmt.Println(pgTag)

	return nil, nil
}
