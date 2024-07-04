package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	db     *pgx.Conn
	DSN    string
	Now    func() time.Time
	ctx    context.Context
	cancel func()
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
	db.db, err = pgx.Connect(db.ctx, db.DSN)
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
		return db.db.Close(db.ctx)
	}
	return nil
}

func (db *DB) Migrate() error {

	migrationTableQuery := `
		CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(100) PRIMARY KEY
		)
	`
	_, err := db.db.Exec(db.ctx, migrationTableQuery)
	if err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	names, err := ReadMigrationDir("migrations", "sql")
	if err != nil {
		return err
	}

	if len(names) == 0 {
		return fmt.Errorf("not sql files found")
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

func (db *DB) migrateFile(name string) error {
	ctx := context.Background()
	tx, err := db.db.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	selectMigration := `SELECT COUNT(*) FROM migrations where name = $1`
	var n int
	err = tx.QueryRow(ctx, selectMigration, name).Scan(&n)
	if err != nil {
		return fmt.Errorf("QueryRow failed: %w", err)
	}

	if n != 0 {
		return nil // migration alreay applied, skip
	}

	buf, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, string(buf))
	if err != nil {
		return err
	}

	insertMigrationQuery := `INSERT INTO migrations(name) VALUES($1)`
	_, err = tx.Exec(ctx, insertMigrationQuery, name)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
