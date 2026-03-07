---
weight: 2
bookFlatSection: true
title: "Commands"
---

# Commands Overview

allbctl is organized into several main commands and subcommands.

## Main Commands

- **`allbctl status`** - Display system information (see [Status Command](../status))
- **`allbctl bootstrap`** - Manage development environment setup (see [Bootstrap Command](../bootstrap))
- **`allbctl version`** - Show version and commit info
- **`allbctl completion`** - Generate shell completion scripts (bash, zsh, fish, PowerShell)
- **`allbctl gen-docs`** - Generate CLI reference documentation

## Status Subcommands

- **`allbctl status runtimes`** - Show detected programming runtimes
- **`allbctl status projects`** - Show git repositories in ~/src
- **`allbctl status list-packages`** - Show package counts
- **`allbctl status db`** - Show detected databases
- **`allbctl status cloud-native`** - Show cloud CLI tools (AWS, GCP, Azure, kubectl)
- **`allbctl status containers`** - Show container runtimes and virtualization
- **`allbctl status git`** - Show git global configuration
- **`allbctl status network`** - Show network interfaces and connectivity
- **`allbctl status ports`** - Show listening TCP/UDP ports
- **`allbctl status security`** - Show SSH keys, GPG keys, and keyring
- **`allbctl status systemctl`** - Show systemd service counts

## Bootstrap Subcommands

- **`allbctl bootstrap status`** - Check bootstrap status
- **`allbctl bootstrap install`** - Install development environment
- **`allbctl bootstrap reset`** - Reset configuration

## Global Flags

- `--config string` - Config file path (default: `$HOME/.allbctl.yaml`)
- `--help, -h` - Show help

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

# Show projects (all repos)
allbctl status projects

# Show packages
allbctl status list-packages

# Show databases
allbctl status db
allbctl status db sqlite3 --detail

# Show network
allbctl status network

# Show containers
allbctl status containers

# Show listening ports
allbctl status ports

# Show version
allbctl version
```
