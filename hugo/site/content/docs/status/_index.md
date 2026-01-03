---
weight: 3
bookFlatSection: true
title: "Status Command"
---

# Status Command

The `allbctl status` command displays comprehensive system information, similar to neofetch but with more development-focused details.

## Basic Usage

```bash
allbctl status
```

## Output Sections

### System Information
- Operating system and version
- Hostname
- Shell
- Terminal
- CPU details (model, cores, architecture)
- GPU information
- Memory usage
- Hardware details

### Development Tools
- **Runtimes**: Detected programming languages (Python, Go, Node.js, Java, etc.)
- **Databases**: Detected database systems and their status
- **Package Managers**: System and programming package managers
- **Packages**: Package counts from all detected managers
- **Projects**: Git repositories in ~/src directory

### Network Information
- Network interfaces
- IP addresses
- Router information
- Connection type

### Additional Information
- Browsers: Detected web browsers with versions
- AI Agents: Detected AI coding assistants

## Subcommands

All status subcommands show specific sections of the main status output:

- [Runtimes](runtimes)
- [Projects](projects)
- [Packages](packages)
- [Databases](databases)
