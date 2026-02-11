package xlogging

import (
	"context"
	"io"
	"log/slog"
)

// Logger is the interface for structured logging.
type Logger interface {
	// Debug logs at debug level.
	Debug(msg string, args ...any)
	// Info logs at info level.
	Info(msg string, args ...any)
	// Warn logs at warn level.
	Warn(msg string, args ...any)
	// Error logs at error level.
	Error(msg string, args ...any)

	// DebugContext logs at debug level with context.
	DebugContext(ctx context.Context, msg string, args ...any)
	// InfoContext logs at info level with context.
	InfoContext(ctx context.Context, msg string, args ...any)
	// WarnContext logs at warn level with context.
	WarnContext(ctx context.Context, msg string, args ...any)
	// ErrorContext logs at error level with context.
	ErrorContext(ctx context.Context, msg string, args ...any)

	// With returns a new Logger with the given attributes.
	With(args ...any) Logger
	// WithGroup returns a new Logger with the given group name.
	WithGroup(name string) Logger
	// Handler returns the underlying slog.Handler.
	Handler() slog.Handler
}

// logger is the concrete implementation of Logger.
type logger struct {
	slog *slog.Logger
}

// New creates a new Logger with the given options.
func New(opts ...Option) Logger {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return newLoggerFromConfig(cfg)
}

// Default creates a new Logger with auto-detected configuration.
// It reads XLOG_ENV and XLOG_LEVEL environment variables.
func Default() Logger {
	return New()
}

// newLoggerFromConfig creates a logger from the given configuration.
func newLoggerFromConfig(cfg *config) Logger {
	handler := createHandler(cfg)
	return &logger{
		slog: slog.New(handler),
	}
}

// createHandler creates the appropriate handler chain based on config.
func createHandler(cfg *config) slog.Handler {
	var baseHandler slog.Handler

	if cfg.shouldUseJSON() {
		baseHandler = slog.NewJSONHandler(cfg.output, &slog.HandlerOptions{
			Level:     cfg.level,
			AddSource: cfg.addSource,
		})
	} else if cfg.shouldUseColor() {
		baseHandler = newColorHandler(cfg.output, &colorHandlerOptions{
			Level:       cfg.level,
			AddSource:   cfg.addSource,
			ContextKeys: cfg.contextKeys,
		})
		// Color handler handles context keys directly, no need to wrap
		return baseHandler
	} else {
		baseHandler = slog.NewTextHandler(cfg.output, &slog.HandlerOptions{
			Level:     cfg.level,
			AddSource: cfg.addSource,
		})
	}

	// Wrap with context handler if context keys are specified
	if len(cfg.contextKeys) > 0 {
		return newContextHandler(baseHandler, cfg.contextKeys)
	}

	return baseHandler
}

// Debug logs at debug level.
func (l *logger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

// Info logs at info level.
func (l *logger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

// Warn logs at warn level.
func (l *logger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

// Error logs at error level.
func (l *logger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

// DebugContext logs at debug level with context.
func (l *logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.slog.DebugContext(ctx, msg, args...)
}

// InfoContext logs at info level with context.
func (l *logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.slog.InfoContext(ctx, msg, args...)
}

// WarnContext logs at warn level with context.
func (l *logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.slog.WarnContext(ctx, msg, args...)
}

// ErrorContext logs at error level with context.
func (l *logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.slog.ErrorContext(ctx, msg, args...)
}

// With returns a new Logger with the given attributes.
func (l *logger) With(args ...any) Logger {
	return &logger{
		slog: l.slog.With(args...),
	}
}

// WithGroup returns a new Logger with the given group name.
func (l *logger) WithGroup(name string) Logger {
	return &logger{
		slog: l.slog.WithGroup(name),
	}
}

// Handler returns the underlying slog.Handler.
func (l *logger) Handler() slog.Handler {
	return l.slog.Handler()
}

// Discard returns a Logger that discards all log output.
func Discard() Logger {
	return &logger{
		slog: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}
