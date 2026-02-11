package xlogging

import (
	"log/slog"
	"strings"
)

// Level represents a logging level.
type Level = slog.Level

// Logging levels.
const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

// ParseLevel parses a level string into a Level.
// Supported values (case-insensitive): debug, info, warn, warning, error.
// Returns LevelInfo if the string is not recognized.
func ParseLevel(s string) Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}
