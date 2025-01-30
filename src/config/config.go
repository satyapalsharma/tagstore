package config

import (
	"fmt"
	"os"
	"strings"

	"tagstore/src/lib"
)

// Config holds application configuration
type Config struct {
	DatabasePath string
	LogLevel     lib.Level
}

// Load loads configuration from environment variables and defaults
func Load() (*Config, error) {
	cfg := &Config{
		DatabasePath: "",
		LogLevel:     lib.LevelInfo,
	}

	// Database path
	if dbPath := os.Getenv("TAGSTORE_DB"); dbPath != "" {
		cfg.DatabasePath = dbPath
	} else {
		defaultPath, err := lib.DefaultDBPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get default database path: %w", err)
		}
		cfg.DatabasePath = defaultPath
	}

	// Log level
	if logLevel := os.Getenv("TAGSTORE_LOG_LEVEL"); logLevel != "" {
		level, err := parseLogLevel(logLevel)
		if err != nil {
			return nil, fmt.Errorf("invalid log level: %w", err)
		}
		cfg.LogLevel = level
	}

	return cfg, nil
}

func parseLogLevel(s string) (lib.Level, error) {
	switch strings.ToLower(s) {
	case "debug":
		return lib.LevelDebug, nil
	case "info":
		return lib.LevelInfo, nil
	case "warn", "warning":
		return lib.LevelWarn, nil
	case "error":
		return lib.LevelError, nil
	default:
		return lib.LevelInfo, fmt.Errorf("unknown log level: %s", s)
	}
}
