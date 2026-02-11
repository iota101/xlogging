# xlogging - Style and Conventions

## Code Style
- Standard Go formatting (`go fmt`)
- Functional options pattern for configuration
- Interface-based design for testability

## API Patterns
- Methods mirror `log/slog` API (Debug, Info, Warn, Error + Context variants)
- `With(args...)` and `WithGroup(name)` return new Logger instances
- Options use `With` prefix: `WithEnv()`, `WithLevel()`, `WithOutput()`
- Context helpers use `With` prefix: `WithRequestID()`, `WithTraceID()`
- Context getters use `Get` prefix: `GetRequestID()`, `GetTraceID()`

## Testing Patterns
- TestLogger captures all log entries for assertions
- `HasEntry(level, msg)` for simple checks
- `HasEntryWithAttr(level, msg, key, value)` for attribute checks
- Child loggers from With()/WithGroup() share entries with parent

## Naming Conventions
- Exported types: `Logger`, `Level`, `Env`, `ContextKey`, `Option`
- Exported constants: `LevelDebug`, `EnvProduction`, `KeyRequestID`
- Private types: `logger`, `config`, `contextHandler`, `colorHandler`
