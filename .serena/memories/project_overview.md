# xlogging - Project Overview

## Purpose
Structured logging library for Go wrapping `log/slog` with environment-aware output formatting and context propagation.

## Tech Stack
- **Language**: Go 1.25
- **Module**: `github.com/iota101/xlogging`
- **Dependencies**: `golang.org/x/term` (terminal detection)

## Project Structure
```
xlogging/
├── logger.go      # Logger interface and implementation
├── options.go     # Functional options (WithEnv, WithLevel, etc.)
├── level.go       # Level type alias and ParseLevel
├── env.go         # Environment detection (XLOG_ENV, XLOG_LEVEL)
├── context.go     # Context keys and helpers (WithRequestID, etc.)
├── handler.go     # contextHandler for extracting context values
├── color.go       # colorHandler for terminal output
├── mock.go        # TestLogger for testing
├── xlogging_test.go
├── Taskfile.yml   # Task commands
├── README.md
└── CLAUDE.md
```

## Key Types
- `Logger` - Main interface for logging
- `Env` - Environment type (production/staging/development)
- `Level` - Log level (alias for slog.Level)
- `ContextKey` - Type for context keys
- `TestLogger` - Logger implementation for testing
- `Option` - Functional option type

## Key Functions
- `New(opts ...Option) Logger` - Create logger with options
- `Default() Logger` - Create from environment variables
- `Discard() Logger` - Silent logger
- `NewTestLogger() *TestLogger` - Create test logger
- `ParseLevel(s string) Level` - Parse level string
- `WithRequestID/TraceID/SpanID/UserID(ctx, value)` - Context helpers

## Environment Variables
| Variable | Values | Default |
|----------|--------|---------|
| `XLOG_ENV` | production/staging/development | development |
| `XLOG_LEVEL` | debug/info/warn/error | env-dependent |

## Output Formats
| Environment | Format | Colors | Default Level |
|-------------|--------|--------|---------------|
| production | JSON | No | Info |
| staging | Text | Yes | Debug |
| development | Text | Yes | Debug |
