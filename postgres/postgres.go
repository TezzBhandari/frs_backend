package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// NullTime represents a helper wrapper for time.Time. It automatically converts
// time fields to/from RFC 3339 format. Also supports NULL for zero time.
type NullTime time.Time

type DB struct {
	db        *pgxpool.Pool
	DSN       string
	Now       func() time.Time
	ctx       context.Context
	cancel    func()
	snowflake *snowflake.Node
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() error {
	var err error

	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	db.snowflake, err = snowflake.NewNode(1)

	if err != nil {
		return err
	}

	db.db, err = pgxpool.New(db.ctx, db.DSN)

	if err != nil {
		return err
	}

	err = db.Migrate()

	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func (db *DB) Close() error {
	defer db.cancel()
	if db.db != nil {
		db.db.Close()
	}
	return nil
}

func (db *DB) Migrate() error {
	migrationTableQuery := `
		CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(100) PRIMARY KEY,
			size INTEGER NOT NULL
		)
	`
	_, err := db.db.Exec(db.ctx, migrationTableQuery)
	if err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	path += "/postgres/migrations"
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return fmt.Errorf("%s path does not exist", path)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	names, err := ReadMigrationDir(path, "sql")
	if err != nil {
		return err
	}

	if len(names) == 0 {
		return fmt.Errorf("no sql files found")
	}

	for _, name := range names {
		err := db.migrateFile(name)
		if err != nil {
			return fmt.Errorf("migration error: name:%q, error: %w", name, err)
		}
	}

	return nil
}

func ReadMigrationDir(dirName, ext string) ([]string, error) {
	var files []string
	dirEntries, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		filename := entry.Name()
		fileParts := strings.Split(filename, ".")
		fileExtension := fileParts[len(fileParts)-1]
		if fileExtension == ext {
			files = append(files, filename)
		}
	}

	return files, nil
}

func (db *DB) migrateFile(filename string) error {
	ctx := context.Background()
	tx, err := db.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	selectMigration := `SELECT name, size, COUNT(name) AS n FROM migrations WHERE name = $1 GROUP BY name;`
	var name string
	var size int
	var n int
	err = tx.QueryRow(ctx, selectMigration, filename).Scan(&name, &size, &n)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			log.Info().Msg(fmt.Sprintf("migration file %s is not applied", filename))
		default:
			return fmt.Errorf("QueryRow failed: %w", err)
		}

	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	path += "/postgres/migrations/" + filename

	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if n == 0 {
		_, err = tx.Exec(ctx, string(buf))
		if err != nil {
			return err
		}

		insertMigrationQuery := `INSERT INTO migrations(name, size) VALUES($1, $2)`
		_, err = tx.Exec(ctx, insertMigrationQuery, filename, len(buf))
		if err != nil {
			return err
		}
	} else if n != 0 && size == len(buf) {
		log.Info().Msg(fmt.Sprintf("migration file %s already applied", name))
		return nil // migration already applied, skip
	} else {
		_, err = tx.Exec(ctx, string(buf))
		if err != nil {
			return err
		}

		insertMigrationQuery := `UPDATE migrations SET size = $1 WHERE  name = $2;`
		_, err = tx.Exec(ctx, insertMigrationQuery, len(buf), filename)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

type Tx struct {
	pgx.Tx
	db  *DB
	Now time.Time
}

func (db *DB) BeginTx(ctx context.Context, txOpts pgx.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, txOpts)

	if err != nil {
		return nil, err
	}

	return &Tx{
		db:  db,
		Tx:  tx,
		Now: db.Now().UTC().Truncate(time.Second),
	}, nil

}
