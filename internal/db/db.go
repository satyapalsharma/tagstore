package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps the database connection by embedding *sql.DB.
type DB struct {
	*sql.DB
}

// InitDB initializes the SQLite database connection and creates the schema.
func InitDB(dbPath string) (*DB, error) {
	if dbPath == "" {
		dbPath = "tags.db"
	}

	// Ensure the directory of dbPath exists.
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return nil, fmt.Errorf("init db: failed to create directory %s: %w", dir, err)
		}
	}

	// Open the SQLite database using modernc.org/sqlite with foreign keys enabled.
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", dbPath)
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}

	// Set max open connections to 1 to avoid SQLite locking issues.
	sqlDB.SetMaxOpenConns(1)

	// Wrap the connection.
	db := &DB{DB: sqlDB}

	// Create tables.
	if err := db.createTables(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("init db: %w", err)
	}

	return db, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

func (d *DB) createTables() error {
	tagTable := `
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	)`

	fileTable := `
	CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT UNIQUE NOT NULL
	)`

	fileTagTable := `
	CREATE TABLE IF NOT EXISTS file_tags (
		file_id INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
		tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
		PRIMARY KEY (file_id, tag_id)
	)`

	for _, stmt := range []string{tagTable, fileTable, fileTagTable} {
		if _, err := d.Exec(stmt); err != nil {
			return fmt.Errorf("create tables: %w", err)
		}
	}

	return nil
}
