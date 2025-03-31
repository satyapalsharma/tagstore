# Contributing to TagStore

Thank you for your interest in contributing to TagStore! This document provides guidelines and instructions for contributing.

## Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/<your-username>/tagstore.git
   cd tagstore
   ```

2. **Install Go** (version 1.21 or later recommended)
   ```bash
   go version
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Build the project**
   ```bash
   go build -o tagstore ./main.go
   ```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./src/service
```

## Code Style

- Follow standard Go formatting: `gofmt` and `goimports`
- Use meaningful variable and function names
- Write comments for exported functions and types
- Keep functions focused and reasonably sized
- Handle errors explicitly; do not ignore them
- Use the project's internal logger (`src/lib/logger.go`) for logging
- Follow the existing package structure

## Pull Request Guidelines

1. **Create a feature branch** from `main`
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes** and ensure:
   - All tests pass (`go test ./...`)
   - Code is formatted (`gofmt -w .`)
   - Lint warnings are addressed (`go vet ./...`)

3. **Commit with a clear message**
   ```bash
   git commit -m "feat: add support for boolean tag queries"
   ```

4. **Push to your fork** and open a Pull Request against `main`

5. **PR Description** should include:
   - What the change does
   - Why it is needed
   - Any breaking changes
   - Screenshots or examples if applicable

6. **Keep PRs focused** — one feature or fix per PR

7. **Be responsive** to review comments and CI feedback
