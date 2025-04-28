package lib

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the database connection
type DB struct {
	db *sql.DB
}

var instance *DB

// Connect initializes the SQLite database connection
func Connect(dbPath string) (*DB, error) {
	var err error
	instance = &DB{}

	instance.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	instance.db.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes
	instance.db.SetMaxIdleConns(1)

	// Enable WAL mode for better concurrency
	if _, err := instance.db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Enable foreign keys
	if _, err := instance.db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := instance.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return instance, nil
}

// GetDB returns the singleton database instance
func GetDB() (*DB, error) {
	if instance == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return instance, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// RawDB returns the underlying *sql.DB instance
func (d *DB) RawDB() *sql.DB {
	return d.db
}

// Exec executes a query without returning rows
func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(query, args...)
}

// Query executes a query that returns rows
func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(query, args...)
}

// Begin starts a transaction
func (d *DB) Begin() (*sql.Tx, error) {
	return d.db.Begin()
}

func (d *DB) createTables() error {
	tagTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	fileTable := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT NOT NULL UNIQUE,
		hash TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	fileTagTable := `
	CREATE TABLE IF NOT EXISTS file_tags (
		file_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (file_id, tag_id),
		FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	)`

	tagIndex := `CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name)`
	fileIndex := `CREATE INDEX IF NOT EXISTS idx_files_path ON files(path)`
	fileTagFileIndex := `CREATE INDEX IF NOT EXISTS idx_file_tags_file_id ON file_tags(file_id)`
	fileTagTagIndex := `CREATE INDEX IF NOT EXISTS idx_file_tags_tag_id ON file_tags(tag_id)`

	for _, stmt := range []string{tagTable, fileTable, fileTagTable, tagIndex, fileIndex, fileTagFileIndex, fileTagTagIndex} {
		if _, err := d.db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

// DefaultDBPath returns the default database path
func DefaultDBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(home, ".tagstore")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(configDir, "tagstore.db"), nil
}
