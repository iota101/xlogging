package xlogging

import (
	"io"
	"os"
)

// config holds the logger configuration.
type config struct {
	env         Env
	level       Level
	output      io.Writer
	contextKeys []ContextKey
	addSource   bool
	useColor    *bool // nil means auto-detect
}

// defaultConfig returns the default configuration.
func defaultConfig() *config {
	env := detectEnv()
	return &config{
		env:         env,
		level:       detectLevel(env),
		output:      os.Stderr,
		contextKeys: nil,
		addSource:   false,
		useColor:    nil,
	}
}

// Option is a functional option for configuring a Logger.
type Option func(*config)

// WithEnv sets the environment for the logger.
// This affects the output format (JSON for production, text for others).
func WithEnv(env Env) Option {
	return func(c *config) {
		c.env = env
	}
}

// WithLevel sets the minimum log level.
func WithLevel(level Level) Option {
	return func(c *config) {
		c.level = level
	}
}

// WithOutput sets the output writer for the logger.
func WithOutput(w io.Writer) Option {
	return func(c *config) {
		c.output = w
	}
}

// WithContextKeys sets the context keys to extract from context.Context.
func WithContextKeys(keys ...ContextKey) Option {
	return func(c *config) {
		c.contextKeys = keys
	}
}

// WithSource enables or disables source code location in log entries.
func WithSource(enabled bool) Option {
	return func(c *config) {
		c.addSource = enabled
	}
}

// WithColor explicitly enables or disables colored output.
// By default, color is auto-detected based on terminal support.
func WithColor(enabled bool) Option {
	return func(c *config) {
		c.useColor = &enabled
	}
}

// shouldUseColor determines if color output should be used.
func (c *config) shouldUseColor() bool {
	if c.useColor != nil {
		return *c.useColor
	}
	// Auto-detect: use color in non-production environments when output is a terminal
	if c.env == EnvProduction {
		return false
	}
	// Check if output is stdout or stderr and is a terminal
	if c.output == os.Stdout {
		return isStdoutTerminal()
	}
	if c.output == os.Stderr {
		return isStderrTerminal()
	}
	return false
}

// shouldUseJSON determines if JSON format should be used.
func (c *config) shouldUseJSON() bool {
	return c.env == EnvProduction
}
