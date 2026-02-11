# xlogging - Task Completion Checklist

## Before Committing
- [ ] `go fmt ./...` - Code is formatted
- [ ] `go vet ./...` - No vet issues
- [ ] `go test ./...` - All tests pass
- [ ] `go test -race ./...` - No race conditions

## Before PR/Release
- [ ] `task check` - All checks pass
- [ ] Tests cover new functionality
- [ ] README.md updated if API changed
- [ ] CLAUDE.md updated if needed

## Quick Check
```bash
task check  # Runs fmt + vet + test:cover
```
