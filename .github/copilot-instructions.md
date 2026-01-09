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
  - Run `go run main.go` as the starting point to allbctl CLI while developing
  - Run `go fmt` to format code
  - Run `go vet` to check for common issues
  - Run `go test` to execute tests
  - Run `go mod tidy` to manage dependencies
  - Follow Go conventions and idioms
  - **BEFORE COMMITTING**: Run `make lint` to catch issues early and prevent CI/CD failures
  - Final checks should be made on built allbctl (Makefile builds to `bin/allbctl`)

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


