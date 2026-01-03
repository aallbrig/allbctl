---
weight: 2
bookFlatSection: true
title: "Commands"
---

# Commands Overview

allbctl is organized into several main commands and subcommands.

## Main Commands

- **`allbctl status`** - Display system information
- **`allbctl bootstrap`** - Manage development environment setup

## Status Subcommands

- **`allbctl status runtimes`** - Show detected programming runtimes
- **`allbctl status projects`** - Show git repositories in ~/src
- **`allbctl status list-packages`** - Show package counts
- **`allbctl status db`** - Show detected databases

## Bootstrap Subcommands

- **`allbctl bootstrap status`** - Check bootstrap status
- **`allbctl bootstrap install`** - Install development environment
- **`allbctl bootstrap reset`** - Reset configuration

## Global Flags

- `--config string` - Config file path (default: $HOME/.allbctl.yaml)
- `--help, -h` - Show help
- `--version, -v` - Show version

## Quick Reference

```bash
# System info
allbctl status

# Dev environment setup
allbctl bootstrap install

# Check what's installed
allbctl bootstrap status

# Show runtimes
allbctl status runtimes

# Show projects
allbctl status projects

# Show packages
allbctl status list-packages

# Show databases
allbctl status db
allbctl status db sqlite3 --detail
```
