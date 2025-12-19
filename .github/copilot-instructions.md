# Development Expectations

## Coding Process

All code development follows Test-Driven Development (TDD):

1. Write failing tests first
2. Write production code iteratively to pass the tests
3. Refactor as needed while keeping tests passing
4. Ensure all tests pass before completing work

## Go-Specific Requirements

All Go code must:

- Have corresponding tests
- Be formatted using the Go CLI (`go fmt`, `go version go1.23.5 linux/amd64`)
- Follow standard Go practices:
  - Ensure `Makefile` is up to date
  - Run `go run main.go` as the entrypoint to allbctl CLI while developing
  - Run `go fmt` to format code
  - Run `go vet` to check for common issues
  - Run `go test` to execute tests
  - Run `go mod tidy` to manage dependencies
  - Follow Go conventions and idioms
  - Final checks should be on built allbctl (Makefile builds to `bin/allbctl`)

## Documentation

- New functionality must be documented in README.md
- When refactoring or changing functionality, critically review if README.md needs updates
- Focus documentation on observable side effects and behavior visible to users
- Internal implementation changes that don't affect user-facing behavior may not require documentation updates

