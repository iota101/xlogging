# xlogging

[![CI](https://github.com/iota101/xlogging/actions/workflows/ci.yml/badge.svg)](https://github.com/iota101/xlogging/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/iota101/xlogging.svg)](https://pkg.go.dev/github.com/iota101/xlogging)
[![Go Report Card](https://goreportcard.com/badge/github.com/iota101/xlogging)](https://goreportcard.com/report/github.com/iota101/xlogging)
[![Latest Tag](https://img.shields.io/github/v/tag/iota101/xlogging?label=latest&sort=semver)](https://github.com/iota101/xlogging/tags)

Structured logging library for Go wrapping `log/slog` with environment-aware output formatting and context propagation.

## Features

- **Environment-aware formatting**: JSON for production, colored text for development/staging
- **Context propagation**: Automatically extracts request_id, trace_id, span_id, user_id from context
- **Functional options**: Clean configuration with `WithEnv()`, `WithLevel()`, `WithContextKeys()`
- **Testing support**: `TestLogger` captures entries for assertions, `Discard()` for silent logging
- **slog compatible**: Full compatibility with `log/slog` patterns

## Installation

```bash
go get github.com/iota101/xlogging
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/iota101/xlogging"
)

func main() {
    // Auto-detect from XLOG_ENV and XLOG_LEVEL
    log := xlogging.Default()
    log.Info("server started", "port", 8080)

    // With context propagation
    ctx := xlogging.WithRequestID(context.Background(), "req-123")
    log.InfoContext(ctx, "request received")
}
```

## Configuration

### Environment Variables

| Variable | Values | Default |
|----------|--------|---------|
| `XLOG_ENV` | `production`, `staging`, `development` | `development` |
| `XLOG_LEVEL` | `debug`, `info`, `warn`, `error` | Depends on env |

### Default Levels by Environment

| Environment | Format | Colors | Default Level |
|-------------|--------|--------|---------------|
| production | JSON | No | Info |
| staging | Text | Yes | Debug |
| development | Text | Yes | Debug |

### Functional Options

```go
log := xlogging.New(
    xlogging.WithEnv(xlogging.EnvProduction),     // JSON format
    xlogging.WithLevel(xlogging.LevelDebug),      // Minimum level
    xlogging.WithOutput(os.Stdout),               // Output writer
    xlogging.WithSource(true),                    // Include source location
    xlogging.WithColor(true),                     // Force color output
    xlogging.WithContextKeys(                     // Context keys to extract
        xlogging.KeyRequestID,
        xlogging.KeyTraceID,
    ),
)
```

## API Reference

### Types

```go
// Logger interface
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    DebugContext(ctx context.Context, msg string, args ...any)
    InfoContext(ctx context.Context, msg string, args ...any)
    WarnContext(ctx context.Context, msg string, args ...any)
    ErrorContext(ctx context.Context, msg string, args ...any)
    With(args ...any) Logger
    WithGroup(name string) Logger
    Handler() slog.Handler
}

// Environment
type Env string
const (
    EnvProduction  Env = "production"
    EnvStaging     Env = "staging"
    EnvDevelopment Env = "development"
)

// Context keys
const (
    KeyRequestID ContextKey = "request_id"
    KeyTraceID   ContextKey = "trace_id"
    KeySpanID    ContextKey = "span_id"
    KeyUserID    ContextKey = "user_id"
)
```

### Functions

| Function | Returns | Description |
|----------|---------|-------------|
| `New(opts ...Option)` | `Logger` | Creates logger with options |
| `Default()` | `Logger` | Creates logger from env vars |
| `Discard()` | `Logger` | Creates silent logger |
| `NewTestLogger()` | `*TestLogger` | Creates test logger |
| `ParseLevel(s string)` | `Level` | Parses level string |
| `WithRequestID(ctx, id)` | `context.Context` | Adds request ID to context |
| `WithTraceID(ctx, id)` | `context.Context` | Adds trace ID to context |
| `WithSpanID(ctx, id)` | `context.Context` | Adds span ID to context |
| `WithUserID(ctx, id)` | `context.Context` | Adds user ID to context |

## HTTP Middleware Example

```go
func LoggingMiddleware(log xlogging.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            requestID := uuid.New().String()
            ctx := xlogging.WithRequestID(r.Context(), requestID)

            rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            next.ServeHTTP(rw, r.WithContext(ctx))

            log.InfoContext(ctx, "http",
                slog.Group("request",
                    slog.String("method", r.Method),
                    slog.String("path", r.URL.Path),
                ),
                slog.Group("response",
                    slog.Int("status", rw.statusCode),
                ),
                slog.Duration("duration", time.Since(start)),
            )
        })
    }
}
```

## Testing

```go
func TestHandler(t *testing.T) {
    log := xlogging.NewTestLogger()
    handler := NewHandler(log)

    handler.Process()

    // Check for specific log entries
    if !log.HasEntry(xlogging.LevelInfo, "processed") {
        t.Error("expected log entry")
    }

    // Check entry with attribute
    if !log.HasEntryWithAttr(xlogging.LevelInfo, "processed", "count", 5) {
        t.Error("expected attribute")
    }

    // Get all entries
    entries := log.Entries()

    // Count by level
    infoCount := log.Count(xlogging.LevelInfo)

    // Clear for next test
    log.Clear()
}
```

### TestLogger Methods

| Method | Description |
|--------|-------------|
| `Entries() []LogEntry` | Returns all captured entries |
| `HasEntry(level, msgSubstring)` | Checks if entry exists |
| `HasEntryWithAttr(level, msg, key, value)` | Checks entry with attribute |
| `Count(level) int` | Counts entries by level |
| `Len() int` | Total entry count |
| `Clear()` | Removes all entries |

## Development

```bash
go test ./...              # Run tests
go test -race -cover ./... # Tests with race detector
go fmt ./...               # Format code
go vet ./...               # Check for issues
```

## License

MIT
