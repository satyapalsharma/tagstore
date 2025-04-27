package lib

import (
	"fmt"
	"log"
	"os"
	"sync"
)

// Level represents log severity
type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// Logger provides structured logging
type Logger struct {
	level  Level
	logger *log.Logger
	mu     sync.Mutex
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the default logger
func Init(level Level) {
	once.Do(func() {
		defaultLogger = New(level)
	})
}

// New creates a new Logger instance
func New(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "[tagstore] ", log.LstdFlags),
	}
}


// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.logf(LevelDebug, msg, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.logf(LevelInfo, msg, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.logf(LevelWarn, msg, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= LevelError {
		l.logf(LevelError, msg, args...)
	}
}

func (l *Logger) logf(level Level, msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	prefix := ""
	switch level {
	case LevelDebug:
		prefix = "DEBUG"
	case LevelInfo:
		prefix = "INFO"
	case LevelWarn:
		prefix = "WARN"
	case LevelError:
		prefix = "ERROR"
	}
	l.logger.Printf("%s: %s", prefix, fmt.Sprintf(msg, args...))
}
