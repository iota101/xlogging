# xlogging - Development Commands

## Using Taskfile (recommended)
```bash
task test        # Run tests
task test:cover  # Tests with coverage
task test:verbose # Verbose test output
task check       # All checks (fmt + vet + test)
task lint        # Run golangci-lint
task fmt         # Format code
task vet         # Run go vet
task tidy        # Run go mod tidy
```

## Tagging & Releases
```bash
task tag -- v1.0.0       # Create tag
task tag:push -- v1.0.0  # Create and push tag
task tag:latest          # Show latest tag
```

## Manual Commands
```bash
go test ./...              # Run tests
go test -race -cover ./... # Tests with race detector
go fmt ./...               # Format code
go vet ./...               # Check for issues
go mod tidy                # Tidy dependencies
```
