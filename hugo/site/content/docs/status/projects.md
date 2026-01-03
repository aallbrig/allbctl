---
weight: 2
title: "Projects"
---

# Status Projects

Display git repositories found in the `~/src` directory.

## Usage

```bash
# Show summary (default)
allbctl status projects

# Show all repositories
allbctl status projects --all

# Show only dirty repos (with uncommitted changes)
allbctl status projects --dirty

# Show only clean repos
allbctl status projects --clean
```

## Output

### Summary Mode (Default)
Shows total count and last 5 recently modified repositories:

```
Projects: 5 total (2 dirty)
  Last 5 recently touched:
    ~/src/allbctl*             aallbrig/allbctl             2026-01-02 13:11 EST
    ~/src/stock-market-words*  aallbrig/stock-market-words  2026-01-01 12:40 EST
    ~/src/dice-gnome-redux     aallbrig/dice-gnome-redux    2025-12-21 11:50 EST
    ~/src/godot-mcp            Coding-Solo/godot-mcp        2025-12-20 11:19 EST
    ~/src/dotfiles             aallbrig/dotfiles            2025-12-16 21:18 EST
```

Repositories with uncommitted changes are marked with an asterisk (*).

### Detailed Modes

**All repositories** (`--all`):
```
Total repos: 5

  ~/src/allbctl*             aallbrig/allbctl             2026-01-02 13:11 EST
  ~/src/stock-market-words*  aallbrig/stock-market-words  2026-01-01 12:40 EST
  ~/src/dice-gnome-redux     aallbrig/dice-gnome-redux    2025-12-21 11:50 EST
  ...
```

**Dirty only** (`--dirty`):
```
Total repos: 2

  ~/src/allbctl*             aallbrig/allbctl             2026-01-02 13:11 EST
  ~/src/stock-market-words*  aallbrig/stock-market-words  2026-01-01 12:40 EST
```

**Clean only** (`--clean`):
```
Total repos: 3

  ~/src/dice-gnome-redux     aallbrig/dice-gnome-redux    2025-12-21 11:50 EST
  ~/src/godot-mcp            Coding-Solo/godot-mcp        2025-12-20 11:19 EST
  ~/src/dotfiles             aallbrig/dotfiles            2025-12-16 21:18 EST
```

## Features

- Automatically finds all git repositories in ~/src
- Shows repository path, origin remote, and last modification time
- Marks dirty repositories (uncommitted changes) with *
- Sorts by modification time (most recent first)
- Supports filtering by dirty/clean status

## Integration

The summary output is shown in the "Projects:" section of the main `allbctl status` command.
