# Error Handling Scheme

This document describes the error handling patterns used in the allbctl codebase.

## General Principles

1. **Always check errors**: All error returns should be checked and handled appropriately
2. **Fail fast**: For critical errors, fail immediately with clear error messages
3. **User-friendly messages**: Provide context and actionable information in error messages
4. **Log for debugging**: Use logging for errors that don't need to halt execution

## Error Handling Patterns

### Pattern 1: Fatal Errors (CLI Commands)

For errors in CLI command execution that should stop the program:

```go
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}
```

**When to use**: Main command execution paths, critical operations that cannot continue

**Example**: `cmd/root.go` Execute function

### Pattern 2: Logged Errors (Non-critical failures)

For errors that should be logged but don't stop execution:

```go
if err != nil {
    log.Printf("Warning: %v\n", err)
    // Continue with degraded functionality
}
```

**When to use**: Optional features, best-effort operations, information gathering

**Example**: System information collection in `cmd/status.go`

### Pattern 3: Returned Errors (Library Functions)

For library/utility functions that should propagate errors to callers:

```go
func doSomething() error {
    if err := operation(); err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}
```

**When to use**: Reusable functions, library code, operations that callers need to handle

**Example**: Most functions in `pkg/` directories

### Pattern 4: Intentionally Ignored Errors

For errors that are safe to ignore with justification:

```go
_ = cmd.Help() //nolint:errcheck // Help display never fails in practice
```

**When to use**:
- Help text display (never fails)
- Test setup operations (test will fail if setup is incorrect)
- Optional cleanup operations
- Write operations to buffers (only fails on OOM)

**Requirements**:
- Must include `//nolint:errcheck` comment
- Must include explanation of why it's safe to ignore

### Pattern 5: Test Setup Errors

For test setup that doesn't need explicit error checking:

```go
_ = os.MkdirAll(path, 0755) //nolint:errcheck // Test setup - test will fail if directory missing
```

**When to use**: Test preparation where subsequent test assertions will catch problems

### Pattern 6: Skip Tests on Environment Issues

For tests that depend on environment state:

```go
home, err := os.UserHomeDir()
if err != nil {
    t.Skip("Unable to get user home directory")
}
```

**When to use**: Tests that require specific environment conditions

## Linting Compliance

### errcheck Linter

The `errcheck` linter ensures all error returns are handled. To suppress false positives:

1. **Explicitly assign to `_`**: Shows the error was considered
2. **Add `//nolint:errcheck` comment**: Explains why it's safe
3. **Provide justification**: Comment explains the reasoning

### unused Linter

Functions and code flagged as unused should be:

1. **Removed** if truly unused and not part of planned features
2. **Marked with `//nolint:unused` if**:
   - Part of a public API that will be used in the future
   - Required for interface implementation
   - Used in a way the linter doesn't detect (e.g., reflection)

## Current Issues and Fixes

### cmd/bootstrap.go and cmd/root.go

**Issue**: `cmd.Help()` error not checked

**Fix**: These are placeholder commands that just display help. Help display never fails in practice, so we use:
```go
_ = cmd.Help() //nolint:errcheck // Help always succeeds
```

### cmd/projects_test.go

**Issue**: `os.MkdirAll()` errors not checked in test setup

**Fix**: Test setup errors don't need explicit checking because:
1. If directory creation fails, subsequent test operations will fail
2. Tests run in isolated environments where filesystem operations should succeed
3. Explicit error handling adds noise without value

```go
_ = os.MkdirAll(repo1, 0755) //nolint:errcheck // Test setup
```

**Issue**: `os.UserHomeDir()` errors not checked

**Fix**: Tests that depend on home directory should skip if unavailable:
```go
home, err := os.UserHomeDir()
if err != nil {
    t.Skip("Unable to get user home directory")
}
```

### cmd/status.go

**Issue**: `formatUptime()` and `printPackageCountInline()` functions are unused

**Analysis**: These functions were likely planned features that were never integrated:
- `formatUptime()`: Formats duration for display but is never called
- `printPackageCountInline()`: Formats package counts but superseded by different implementation

**Options**:
1. Remove the functions (clean up dead code)
2. Keep with `//nolint:unused` if they're part of planned features
3. Integrate them into the active codebase if they add value

**Recommendation**: Since the PR comment asks to "undo these changes and state why", these functions should be kept and marked as planned features OR the person requesting the review should clarify if they want them integrated or removed differently.

## Summary

This error handling scheme balances:
- **Safety**: Critical errors stop execution with clear messages
- **Usability**: Non-critical failures degrade gracefully
- **Code quality**: Linting tools ensure errors are explicitly considered
- **Pragmatism**: Unnecessary error handling is avoided with clear justification
