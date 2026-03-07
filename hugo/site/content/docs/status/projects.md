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

## Integration

The `allbctl status` command includes a projects section limited to 5 repos:

```
Projects: 12 total (7 dirty)
  Recently touched (5):
    ~/src/allbctl*  ...
    ...
```

To see all repositories without the limit, run `allbctl status projects` directly.
