---
title: "allbctl Documentation"
type: docs
---

# allbctl Documentation

**allbctl** (aka allbrightctl) is a CLI for managing Andrew Allbright's development environment across multiple platforms.

## Quick Start

```bash
# Check system status
allbctl status

# Bootstrap development environment
allbctl bootstrap install

# Show detected runtimes
allbctl status runtimes

# Show git projects
allbctl status projects

# Show installed packages
allbctl status list-packages

# Show detected databases
allbctl status db
```

## Features

- **System Information**: Neofetch-style system info display
- **Bootstrap**: Automated development environment setup
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Package Detection**: Detects packages from multiple package managers
- **Runtime Detection**: Finds installed programming languages and tools
- **Database Detection**: Discovers installed database systems
- **Project Management**: Tracks git repositories

## Platform Support

- Linux (Ubuntu, Arch, Debian, Fedora, etc.)
- macOS
- Windows 10/11

## Installation

See [Installation Guide](/docs/getting-started/installation) for details.
