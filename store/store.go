// Package store keeps scraped ManagerZone data in a local SQLite database.
//
// The schema is meant to be queried directly too (mz sql, sqlite3): times are
// UTC text in RFC3339 ("2026-10-12T10:53:00Z"), so SQLite's date functions and
// plain string comparison work on them.
package store

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaSQL string

// schemaVersion is stored in PRAGMA user_version. Bump it with a migration
// step in migrate when schema.sql changes.
const schemaVersion = 1

type DB struct {
	*sql.DB
}

// DefaultPath is $MZ_DB, else $XDG_DATA_HOME/mz/mz.db, else ~/.local/share/mz/mz.db.
func DefaultPath() (string, error) {
	if p := os.Getenv("MZ_DB"); p != "" {
		return p, nil
	}
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dir, "mz", "mz.db"), nil
}

// Open opens (creating if needed) the database at path and brings its schema up to date.
func Open(path string) (*DB, error) {
	return open(path, false)
}

// OpenReadOnly opens an existing database; statements that write fail.
func OpenReadOnly(path string) (*DB, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("no database at %s (run 'mz forum sync' first): %w", path, err)
	}
	return open(path, true)
}

func open(path string, readOnly bool) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	dsn := "file:" + path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	if readOnly {
		dsn += "&_pragma=query_only(1)"
	}
	sdb, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// One connection: pragmas apply per connection, and SQLite serializes writes anyway.
	sdb.SetMaxOpenConns(1)
	db := &DB{sdb}
	if !readOnly {
		if err := db.migrate(); err != nil {
			sdb.Close()
			return nil, fmt.Errorf("%s: %w", path, err)
		}
	}
	return db, nil
}

func (db *DB) migrate() error {
	var v int
	if err := db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		return err
	}
	switch {
	case v == schemaVersion:
		return nil
	case v > schemaVersion:
		return fmt.Errorf("database schema version %d is newer than this mz (%d)", v, schemaVersion)
	case v == 0:
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		if _, err := tx.Exec(schemaSQL); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
		if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d", schemaVersion)); err != nil {
			return err
		}
		return tx.Commit()
	}
	return fmt.Errorf("no migration from schema version %d", v)
}

const timeFormat = "2006-01-02T15:04:05Z"

func fmtTime(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t.UTC().Format(timeFormat)
}

func parseTime(s sql.NullString) time.Time {
	if !s.Valid {
		return time.Time{}
	}
	t, err := time.Parse(timeFormat, s.String)
	if err != nil {
		return time.Time{}
	}
	return t.In(time.Local)
}

func nullInt(v int) any {
	if v == 0 {
		return nil
	}
	return v
}

func nullStr(s string) any {
	if s == "" {
		return nil
	}
	return s
}
