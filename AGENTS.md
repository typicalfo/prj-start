# AGENTS.md - Development Guidelines

## Build/Test Commands
- **Build**: `make build` (or `go build -o main .` in individual project directories)
- **Run**: `make run` (or `go run .` in individual project directories)  
- **Dev with hot reload**: `make dev` (uses Air for automatic reloading)
- **Test single file**: `go test -v ./path/to/package -run TestFunctionName`
- **Test all**: `make test` (or `go test ./...` in individual project directories)
- **Clean**: `make clean` (removes build artifacts)
- **Kill processes**: `make kill` (kills running servers)
- **Dependencies**: `make deps` (or `go mod tidy && go mod download`)
- **Help**: `make help` (shows all available targets)

## Code Style Guidelines

### Imports
- Group imports: stdlib, third-party, local packages
- Use absolute imports for local packages (e.g., "app/server/domain")
- Avoid unused imports

### Naming Conventions  
- **Packages**: lowercase, single word when possible
- **Functions**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase, descriptive names
- **Constants**: UPPER_SNAKE_CASE for exported
- **Interfaces**: Usually ends with "er" suffix (e.g., `Service`, `Repository`)

### Error Handling
- Always handle errors explicitly
- Use structured logging with colored output (e.g., `github.com/sirupsen/logrus` or `github.com/fatih/color`)
- Return errors from functions, don't panic
- Wrap errors with context using `errors.Wrap` when available

### Testing
- Use testify for assertions and mocking
- Test file naming: `filename_test.go`
- Use table-driven tests for multiple scenarios
- Mock external dependencies

### Project Structure
- Each subdirectory in dev-docs/ is a standalone Go module
- Follow clean architecture patterns (handlers, services, repositories)
- Use Fiber v2 as the web framework
- Separate concerns: handlers, services, repositories, entities

### General
- Use Go 1.23+ features where appropriate
- Prefer composition over inheritance
- Keep functions small and focused
- Add comments for exported functions explaining purpose