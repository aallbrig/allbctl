---
weight: 2
title: "Projects"
---

# Status Projects

Display git repositories found in the `~/src` directory.

## Usage

```bash
# Show all repositories (default)
allbctl status projects

# Limit the number of repositories shown
allbctl status projects --limit 5

# Show all repositories explicitly
allbctl status projects --all

# Show only dirty repos (with uncommitted changes)
allbctl status projects --dirty

# Show only clean repos
allbctl status projects --clean
```

## Flags

| Flag | Description |
|------|-------------|
| `--limit int` | Limit the number of projects shown (0 = no limit, show all) |
| `--all` | Show all detected git repos |
| `--dirty` | Show only repos with uncommitted changes |
| `--clean` | Show only repos with no uncommitted changes |
| `-v, --verbose` | Show detailed information including changed files, CI status, and language breakdown |
| `--languages` | Show language breakdown for each repo (default `true`; use `--languages=false` to hide) |

## Output

### Default Mode
Shows all repositories sorted by most recently modified:

```
Projects: 12 total (7 dirty)
  Recently touched (12):
    ~/src/allbctl*                  aallbrig/allbctl                   2026-03-07 16:19 EST
    ~/src/stock-market-words        aallbrig/stock-market-words        2026-03-07 11:40 EST
    ~/src/dice-gnome-redux          aallbrig/dice-gnome-redux          2025-12-21 11:50 EST
    ~/src/godot-mcp                 Coding-Solo/godot-mcp              2025-12-20 11:19 EST
    ~/src/dotfiles                  aallbrig/dotfiles                  2025-12-16 21:18 EST
    ...
```

Repositories with uncommitted changes are marked with an asterisk (*).

### Verbose Mode (`-v`)

Shows language breakdown, CI check details, and changed files for each repo:

```bash
allbctl status projects -v
```

```
  ~/src/allbctl*  aallbrig/allbctl  2026-03-29 11:09 EDT -0400  [uncommitted changes]  ✓
      Languages: Go: 330.8 KB (74%) | Markdown: 83.7 KB (18%) | YAML: 8.4 KB (1%) | Shell: 7.8 KB (1%)
      ✓ test (ubuntu-latest)
      ✓ lint
      On branch main
      Changes not staged for commit:
      	modified:   cmd/projects.go
```

### Languages (`--languages`)

The `--languages` flag can be used standalone or combined with any filter flag
to show language breakdown without the full verbose output:

```bash
# Languages for all repos
allbctl status projects --all --languages

# Languages for dirty repos only
allbctl status projects --dirty --languages

# Verbose output WITHOUT languages
allbctl status projects -v --languages=false
```

Language detection analyzes tracked files in each repository (excluding vendored directories
like `vendor/` and `node_modules/`). Results are cached per commit SHA in
`~/.cache/allbctl/languages/` so repeated runs are fast.

### With `--limit N`
Shows at most N repositories:

```bash
allbctl status projects --limit 5
```

```
Projects: 12 total (7 dirty)
  Recently touched (5):
    ~/src/allbctl*                  aallbrig/allbctl                   2026-03-07 16:19 EST
    ~/src/stock-market-words        aallbrig/stock-market-words        2026-03-07 11:40 EST
    ~/src/dice-gnome-redux          aallbrig/dice-gnome-redux          2025-12-21 11:50 EST
    ~/src/godot-mcp                 Coding-Solo/godot-mcp              2025-12-20 11:19 EST
    ~/src/dotfiles                  aallbrig/dotfiles                  2025-12-16 21:18 EST
```

### Dirty only (`--dirty`)
```
Projects: 7 total (7 dirty)

  ~/src/allbctl*             aallbrig/allbctl             2026-03-07 16:19 EST
  ~/src/dotfiles*            aallbrig/dotfiles            2025-12-16 21:18 EST
  ...
```

### Clean only (`--clean`)
```
Projects: 5 total

  ~/src/stock-market-words   aallbrig/stock-market-words  2026-03-07 11:40 EST
  ~/src/godot-mcp            Coding-Solo/godot-mcp        2025-12-20 11:19 EST
  ...
```

## Features

- Automatically finds all git repositories in ~/src
- Shows repository path, origin remote, and last modification time
- Marks dirty repositories (uncommitted changes) with *
- Sorts by modification time (most recent first)
- Supports filtering by dirty/clean status
- Verbose mode shows language breakdown, CI check status, and changed files
- Language detection results are cached per-commit for fast subsequent runs
- Cache stored in OS-appropriate location (`~/.cache/allbctl/` on Linux, `~/Library/Caches/allbctl/` on macOS, `%LOCALAPPDATA%\allbctl\` on Windows)

## Integration

The `allbctl status` command includes a projects section limited to 5 repos:

```
Projects: 12 total (7 dirty)
  Recently touched (5):
    ~/src/allbctl*  ...
    ...
```

To see all repositories without the limit, run `allbctl status projects` directly.
