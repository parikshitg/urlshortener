package logger

import (
	"fmt"
	"os"
	"time"
)

// Logger provides basic structured logging
type Logger struct {
	level  string
	format string
}

// New creates a logger
func New(level, format string) *Logger {
	return &Logger{
		level:  level,
		format: format,
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.shouldLog("info") {
		l.log("INFO", msg, fields...)
	}
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...interface{}) {
	if l.shouldLog("error") {
		l.log("ERROR", msg, fields...)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.shouldLog("debug") {
		l.log("DEBUG", msg, fields...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...interface{}) {
	if l.shouldLog("warn") {
		l.log("WARN", msg, fields...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log("FATAL", msg, fields...)
	os.Exit(1)
}

// WithField creates a logger with a field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return l
}

// WithError creates a logger with an error
func (l *Logger) WithError(err error) *Logger {
	return l
}

// WithFields creates a logger with fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return l
}

func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
		"fatal": 4,
	}

	currentLevel := levels[l.level]
	msgLevel := levels[level]
	return msgLevel >= currentLevel
}

func (l *Logger) log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if l.format == "json" {
		// Simple JSON format
		fmt.Printf(`{"timestamp":"%s","level":"%s","msg":"%s"}`, timestamp, level, msg)
		if len(fields) > 0 {
			fmt.Printf(",\"fields\":%v", fields)
		}
		fmt.Println()
	} else {
		// Simple text format
		fmt.Printf("[%s] %s: %s", timestamp, level, msg)
		if len(fields) > 0 {
			fmt.Printf(" %v", fields)
		}
		fmt.Println()
	}
}

// LoadSimpleConfig loads simple logging configuration
func LoadSimpleConfig() (string, string) {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		format = "text"
	}

	return level, format
}
