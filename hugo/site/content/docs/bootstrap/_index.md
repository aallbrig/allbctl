---
weight: 4
bookFlatSection: true
title: "Bootstrap Command"
---

# Bootstrap Command

The `allbctl bootstrap` command automates the setup of a development environment.

## Overview

Bootstrap provides three subcommands for managing your development environment:

- [Status](status) - Check what's installed and what's missing
- [Install](install) - Install and configure development environment  
- [Reset](reset) - Reset bootstrap configuration (removes installed items)

## What Bootstrap Does

### Directory Setup
- Creates `~/src` directory for source code

### Tool Installation
- **git** - Version control system (cross-platform)
- **gh** - GitHub CLI (cross-platform)

### SSH Keys (Optional)
- Generates SSH keys if missing
- Registers keys with GitHub
- **Requires `--register-ssh-keys` flag**

### Dotfiles
- Clones dotfiles repository from GitHub
- Runs installation script (if available)

### Shell Tools
- Detects tools referenced in shell config files
- Shows which are installed vs missing

## Quick Reference

```bash
# Check status
allbctl bootstrap status

# Install everything (except SSH keys)
allbctl bootstrap install

# Install with SSH keys
allbctl bootstrap install --register-ssh-keys

# Reset/remove everything
allbctl bootstrap reset
```

## Platform Support

- **Linux**: apt, dnf, yum, pacman, zypper, apk
- **macOS**: homebrew
- **Windows**: winget, choco, scoop

## Aliases

- `allbctl bs` - Short for `allbctl bootstrap`

