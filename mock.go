package xlogging

import (
	"context"
	"log/slog"
	"strings"
	"sync"
)

// LogEntry represents a captured log entry for testing.
type LogEntry struct {
	Level   Level
	Message string
	Attrs   map[string]any
}

// TestLogger is a Logger implementation for testing that captures log entries.
type TestLogger struct {
	mu      *sync.Mutex
	entries *[]LogEntry
	attrs   map[string]any
	group   string
}

// NewTestLogger creates a new TestLogger for testing.
func NewTestLogger() *TestLogger {
	entries := make([]LogEntry, 0)
	return &TestLogger{
		mu:      &sync.Mutex{},
		entries: &entries,
		attrs:   make(map[string]any),
	}
}

// log adds an entry to the captured logs.
func (t *TestLogger) log(level Level, msg string, args ...any) {
	t.mu.Lock()
	defer t.mu.Unlock()

	entry := LogEntry{
		Level:   level,
		Message: msg,
		Attrs:   make(map[string]any),
	}

	// Copy pre-set attributes
	for k, v := range t.attrs {
		entry.Attrs[k] = v
	}

	// Parse args as key-value pairs
	for i := 0; i < len(args)-1; i += 2 {
		if key, ok := args[i].(string); ok {
			attrKey := key
			if t.group != "" {
				attrKey = t.group + "." + key
			}
			entry.Attrs[attrKey] = args[i+1]
		}
	}

	*t.entries = append(*t.entries, entry)
}

// Debug logs at debug level.
func (t *TestLogger) Debug(msg string, args ...any) {
	t.log(LevelDebug, msg, args...)
}

// Info logs at info level.
func (t *TestLogger) Info(msg string, args ...any) {
	t.log(LevelInfo, msg, args...)
}

// Warn logs at warn level.
func (t *TestLogger) Warn(msg string, args ...any) {
	t.log(LevelWarn, msg, args...)
}

// Error logs at error level.
func (t *TestLogger) Error(msg string, args ...any) {
	t.log(LevelError, msg, args...)
}

// DebugContext logs at debug level with context.
func (t *TestLogger) DebugContext(_ context.Context, msg string, args ...any) {
	t.log(LevelDebug, msg, args...)
}

// InfoContext logs at info level with context.
func (t *TestLogger) InfoContext(_ context.Context, msg string, args ...any) {
	t.log(LevelInfo, msg, args...)
}

// WarnContext logs at warn level with context.
func (t *TestLogger) WarnContext(_ context.Context, msg string, args ...any) {
	t.log(LevelWarn, msg, args...)
}

// ErrorContext logs at error level with context.
func (t *TestLogger) ErrorContext(_ context.Context, msg string, args ...any) {
	t.log(LevelError, msg, args...)
}

// With returns a new TestLogger with the given attributes.
func (t *TestLogger) With(args ...any) Logger {
	t.mu.Lock()
	defer t.mu.Unlock()

	newLogger := &TestLogger{
		mu:      t.mu,      // Share the mutex
		entries: t.entries, // Share the entries slice
		attrs:   make(map[string]any),
		group:   t.group,
	}

	// Copy existing attrs
	for k, v := range t.attrs {
		newLogger.attrs[k] = v
	}

	// Add new attrs
	for i := 0; i < len(args)-1; i += 2 {
		if key, ok := args[i].(string); ok {
			newLogger.attrs[key] = args[i+1]
		}
	}

	return newLogger
}

// WithGroup returns a new TestLogger with the given group name.
func (t *TestLogger) WithGroup(name string) Logger {
	t.mu.Lock()
	defer t.mu.Unlock()

	newGroup := name
	if t.group != "" {
		newGroup = t.group + "." + name
	}

	newLogger := &TestLogger{
		mu:      t.mu,      // Share the mutex
		entries: t.entries, // Share the entries slice
		attrs:   make(map[string]any),
		group:   newGroup,
	}

	// Copy existing attrs
	for k, v := range t.attrs {
		newLogger.attrs[k] = v
	}

	return newLogger
}

// Handler returns nil for TestLogger (not backed by slog.Handler).
func (t *TestLogger) Handler() slog.Handler {
	return nil
}

// Entries returns all captured log entries.
func (t *TestLogger) Entries() []LogEntry {
	t.mu.Lock()
	defer t.mu.Unlock()

	result := make([]LogEntry, len(*t.entries))
	copy(result, *t.entries)
	return result
}

// Clear removes all captured log entries.
func (t *TestLogger) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	*t.entries = (*t.entries)[:0]
}

// HasEntry checks if there's an entry with the given level and message substring.
func (t *TestLogger) HasEntry(level Level, msgSubstring string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, e := range *t.entries {
		if e.Level == level && strings.Contains(e.Message, msgSubstring) {
			return true
		}
	}
	return false
}

// HasEntryWithAttr checks if there's an entry with the given level, message, and attribute.
func (t *TestLogger) HasEntryWithAttr(level Level, msgSubstring string, key string, value any) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, e := range *t.entries {
		if e.Level == level && strings.Contains(e.Message, msgSubstring) {
			if v, ok := e.Attrs[key]; ok && v == value {
				return true
			}
		}
	}
	return false
}

// Count returns the number of entries matching the given level.
func (t *TestLogger) Count(level Level) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	count := 0
	for _, e := range *t.entries {
		if e.Level == level {
			count++
		}
	}
	return count
}

// Len returns the total number of captured entries.
func (t *TestLogger) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(*t.entries)
}
