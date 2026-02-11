# xlogging

Structured logging library wrapping `log/slog` with environment-aware formatting and context propagation.

## Quick Start

```bash
go test ./...              # Run tests
go test -race -cover ./... # Tests with race detector
```

## Structure

```
xlogging/
├── logger.go      # Logger interface and implementation
├── options.go     # Functional options (WithEnv, WithLevel, etc.)
├── level.go       # Level type and ParseLevel
├── env.go         # Environment detection (XLOG_ENV, XLOG_LEVEL)
├── context.go     # Context keys and helpers
├── handler.go     # contextHandler for extracting context values
├── color.go       # colorHandler for terminal output
├── mock.go        # TestLogger for testing
└── xlogging_test.go
```

## API

### Types

```go
type Logger interface {
    Debug/Info/Warn/Error(msg string, args ...any)
    DebugContext/InfoContext/WarnContext/ErrorContext(ctx, msg, args)
    With(args ...any) Logger
    WithGroup(name string) Logger
    Handler() slog.Handler
}

type Env string  // EnvProduction, EnvStaging, EnvDevelopment
type Level = slog.Level  // LevelDebug, LevelInfo, LevelWarn, LevelError
type ContextKey string  // KeyRequestID, KeyTraceID, KeySpanID, KeyUserID
```

### Functions

| Function | Returns | Use Case |
|----------|---------|----------|
| `New(opts...)` | `Logger` | Create with options |
| `Default()` | `Logger` | Auto-detect from env |
| `Discard()` | `Logger` | Silent logger |
| `NewTestLogger()` | `*TestLogger` | Testing |
| `ParseLevel(s)` | `Level` | Parse "debug"/"info"/etc |
| `WithRequestID(ctx, id)` | `context.Context` | Add request_id |
| `WithTraceID(ctx, id)` | `context.Context` | Add trace_id |

### Options

| Option | Description |
|--------|-------------|
| `WithEnv(env)` | Set environment (affects format) |
| `WithLevel(level)` | Set minimum level |
| `WithOutput(w)` | Set output writer |
| `WithContextKeys(keys...)` | Keys to extract from context |
| `WithSource(bool)` | Include source location |
| `WithColor(bool)` | Force color on/off |

## Usage Pattern

```go
// Production setup
log := xlogging.New(
    xlogging.WithEnv(xlogging.EnvProduction),
    xlogging.WithContextKeys(xlogging.KeyRequestID, xlogging.KeyTraceID),
)

// With context
ctx := xlogging.WithRequestID(ctx, "req-123")
log.InfoContext(ctx, "request", "method", "GET")

// Testing
testLog := xlogging.NewTestLogger()
// ... use testLog ...
if !testLog.HasEntry(xlogging.LevelInfo, "expected") {
    t.Error("missing log")
}
```

## Environment Variables

| Variable | Values | Default |
|----------|--------|---------|
| `XLOG_ENV` | production/staging/development | development |
| `XLOG_LEVEL` | debug/info/warn/error | env-dependent |

## Notes

- Production uses JSON format, others use colored text
- Context values (request_id, etc.) are automatically added to logs
- TestLogger shares entries between With()/WithGroup() children
- Discard() for benchmarks or when logging should be suppressed
