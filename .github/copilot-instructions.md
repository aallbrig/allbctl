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
- Be formatted using the Go CLI (`go fmt`, `go version go1.26.1 linux/amd64`)
- Follow standard Go practices:
  - Ensure `Makefile` is up to date
  - Run `go run main.go` as the starting point to allbctl CLI while developing
  - Run `go fmt` to format code
  - Run `go vet` to check for common issues
  - Run `go test` to execute tests
  - Run `go mod tidy` to manage dependencies
  - Follow Go conventions and idioms
  - **BEFORE COMMITTING**: Run `make lint` to catch issues early and prevent CI/CD failures
  - **BEFORE COMMITTING**: Run `govulncheck ./...` to check for known vulnerabilities — CI enforces this and will fail if vulns are present
  - Final checks should be made on built allbctl (Makefile builds to `bin/allbctl`)
  - **AFTER MAKING CHANGES**: Run `make install` to install the latest binary (with version/commit info) to `$GOPATH/bin/allbctl` so it is immediately available to run

## Documentation

### README.md Updates
- New functionality must be documented in README.md
- When refactoring or changing functionality, critically review if README.md needs updates
- Focus documentation on observable side effects and behavior visible to users
- Internal implementation changes that don't affect user-facing behavior may not require documentation updates

### Hugo Site Documentation Updates
When allbctl functionality changes, update the Hugo documentation site:

**Location**: `hugo/site/content/docs/`

**What to update**:
1. **New commands/subcommands**: Create new markdown files in appropriate section
2. **Changed behavior**: Update existing documentation pages
3. **New flags/options**: Add to relevant command documentation
4. **Examples**: Update code examples to reflect current behavior

**Documentation structure**:
```
hugo/site/content/docs/
├── getting-started/
│   └── installation.md        # Installation instructions
├── commands/
│   └── _index.md             # Commands overview
├── status/
│   ├── _index.md             # Status command overview
│   ├── runtimes.md           # status runtimes subcommand
│   ├── projects.md           # status projects subcommand
│   ├── packages.md           # status list-packages subcommand
│   └── databases.md          # status db subcommand
└── bootstrap/
    ├── _index.md             # Bootstrap command overview
    ├── status.md             # bootstrap status subcommand
    ├── install.md            # bootstrap install subcommand
    └── reset.md              # bootstrap reset subcommand
```

**Testing documentation**:
```bash
# Test Hugo site builds
cd hugo/site
hugo --minify

# Preview locally
hugo server --buildDrafts
```

**Deployment**:
- Documentation deploys automatically on new releases via GitHub Actions
- Or manually trigger: Go to Actions → "Deploy Hugo Site to GitHub Pages" → Run workflow
- Verify at: https://aallbrig.github.io/allbctl/

### CLI Reference Auto-generation

The `hugo/site/content/docs/reference/` directory is **auto-generated** from the Cobra command tree — do not edit those files manually.

```bash
# Regenerate after adding/changing commands or flags
make gen-docs

# CI will fail if generated docs are out of sync with code
make check-docs
```

## Windows Cross-Platform Validation

A Vagrant-managed Windows 10 VM is available for smoke testing Windows compatibility.

**VM name**: `windows10` (defined in `Vagrantfile`)

**When to use**: Before cutting a release, or when changing OS-specific code (network detection, ping, paths, shell detection, etc.)

```bash
# Start the VM and run the smoke test
vagrant up windows10
vagrant provision windows10 --provision-with smoke-test

# If binary on the VM is stale, force-copy the new one first
make build-windows  # produces allbctl_windows_amd64.exe in project root
vagrant winrm windows10 -s powershell -c \
  "Copy-Item 'C:\vagrant\allbctl_windows_amd64.exe' 'C:\allbctl-test\allbctl.exe'"
vagrant provision windows10 --provision-with smoke-test

# Halt when done
vagrant halt windows10
```

**Expected result**: `Results: 6 passed, 0 failed` with `Internet: ✓ Connected` in the status output.

**Note on internet check**: Uses `net.DialTimeout("tcp", "8.8.8.8:53", 2s)` — pure Go, no `ping` binary required. This was intentional to avoid cross-OS exec differences.
