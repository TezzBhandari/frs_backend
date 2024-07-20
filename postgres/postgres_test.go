package postgres_test

import (
	"reflect"
	"testing"

	p "github.com/TezzBhandari/frs/postgres"
)

func TestReadMigrationDir(t *testing.T) {
	expected := []string{"user.sql"}
	got, err := p.ReadMigrationDir("migrations", "sql")
	if err != nil {
		t.Errorf("got: %q, want: %q, error: %q", got, expected, err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got: %q, want: %q", got, expected)
	}
}
