package xlogging

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// Env represents the application environment.
type Env string

// Environment constants.
const (
	EnvProduction  Env = "production"
	EnvStaging     Env = "staging"
	EnvDevelopment Env = "development"
)

// Environment variable names.
const (
	envKeyEnv   = "XLOG_ENV"
	envKeyLevel = "XLOG_LEVEL"
)

// detectEnv reads XLOG_ENV and returns the corresponding Env.
// Defaults to EnvDevelopment if not set or unrecognized.
func detectEnv() Env {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(envKeyEnv)))
	switch val {
	case "production", "prod":
		return EnvProduction
	case "staging", "stage":
		return EnvStaging
	case "development", "dev":
		return EnvDevelopment
	default:
		return EnvDevelopment
	}
}

// detectLevel reads XLOG_LEVEL and returns the corresponding Level.
// If not set, returns the default level for the given environment.
func detectLevel(env Env) Level {
	if val := os.Getenv(envKeyLevel); val != "" {
		return ParseLevel(val)
	}
	return defaultLevelForEnv(env)
}

// defaultLevelForEnv returns the default log level for the given environment.
func defaultLevelForEnv(env Env) Level {
	switch env {
	case EnvProduction:
		return LevelInfo
	case EnvStaging, EnvDevelopment:
		return LevelDebug
	default:
		return LevelDebug
	}
}

// isTerminal checks if the given file descriptor is a terminal.
func isTerminal(fd int) bool {
	return term.IsTerminal(fd)
}

// isStdoutTerminal checks if stdout is a terminal.
func isStdoutTerminal() bool {
	return isTerminal(int(os.Stdout.Fd()))
}

// isStderrTerminal checks if stderr is a terminal.
func isStderrTerminal() bool {
	return isTerminal(int(os.Stderr.Fd()))
}
