---
weight: 2
title: "Bootstrap Install"
---

# Bootstrap Install

Install and configure your development environment automatically.

## Usage

```bash
# Basic install (no SSH keys)
allbctl bootstrap install

# With SSH key generation and GitHub registration
allbctl bootstrap install --register-ssh-keys
```

## What Gets Installed

### 1. Directory Creation
Creates `~/src` directory for source code projects.

### 2. Tool Installation
Installs essential development tools:
- **git** - Version control system
- **gh** - GitHub CLI

Installation method varies by platform:
- Linux: Uses system package manager (apt, dnf, pacman, etc.)
- macOS: Uses homebrew
- Windows: Uses winget, choco, or scoop

### 3. SSH Keys (Optional)
Only with `--register-ssh-keys` flag:
- Generates SSH key pair if not exists
- Registers public key with GitHub via gh CLI

### 4. Dotfiles
- Clones dotfiles repository from GitHub
- Default: `https://github.com/aallbrig/dotfiles`
- Location: `~/src/dotfiles`
- Runs `./fresh.sh` installation script if present

## Flags

### --register-ssh-keys
```bash
allbctl bootstrap install --register-ssh-keys
```

Enables SSH key generation and GitHub registration.

**Requirements**:
- GitHub CLI (`gh`) must be installed
- Must authenticate first: `gh auth login`

**What it does**:
1. Checks if SSH key exists at `~/.ssh/id_rsa.pub`
2. Generates key if missing
3. Checks if key is registered with GitHub
4. Registers key if not already registered

**Why it's optional**:
- Not everyone uses SSH for git
- Safer for CI/CD and test environments
- Requires GitHub authentication
- Some prefer manual SSH key management

## Idempotency

Bootstrap install is **fully idempotent** - safe to run multiple times.

### What Happens on Second Run

**Directories**: 
```
✅ Directory already exists - skipped
```

**Tools**:
```
✅ git already installed - skipped
✅ gh already installed - skipped
```

**Dotfiles**:
```
✅ Dotfiles already cloned - skipped
```

**SSH Keys** (with flag):
```
✅ SSH key already registered - skipped
```

The install script will skip any component that's already set up correctly.

## Examples

### First Time Setup
```bash
# Check what's missing
allbctl bootstrap status

# Install everything (except SSH keys)
allbctl bootstrap install

# Verify installation
allbctl bootstrap status
```

Output during install:
```
Applying configuration: Expected Directories
PRESENT /home/user/src
✅ Directory already exists

Applying configuration: Required Tools
INSTALLED git
✅ git already installed

NOT FOUND gh
Installing gh...
✅ gh installed successfully

Applying configuration: Dotfiles
NOT CLONED /home/user/src/dotfiles
Cloning dotfiles...
✅ Dotfiles cloned successfully
Running install script...
✅ Install script completed
```

### With SSH Keys
```bash
# First authenticate with GitHub
gh auth login

# Install with SSH key registration
allbctl bootstrap install --register-ssh-keys
```

Output:
```
...
Applying configuration: SSH Configuration
SSH KEY NOT FOUND
Generating SSH key...
✅ SSH key generated

Registering SSH key with GitHub...
✅ SSH key registered with GitHub
```

### Test Idempotency
```bash
# Run multiple times - should be safe
allbctl bootstrap install
allbctl bootstrap install
allbctl bootstrap install

# All runs after the first will skip existing components
```

## Platform-Specific Notes

### Linux
- Auto-detects package manager (apt, dnf, yum, pacman, etc.)
- May require sudo for package installation
- Full SSH key support

### macOS
- Requires homebrew installed
- Installs via `brew install git gh`
- Full SSH key support

### Windows
- Tries winget, then choco, then scoop
- Some package names may vary
- Dotfiles script requires Git Bash or WSL
- SSH key support via Git Bash or WSL

## Troubleshooting

### gh CLI Installation Fails
```
❌ Failed to install gh
```

**Solutions**:
- Check package manager is working
- Install manually: https://cli.github.com/manual/installation
- Package name may vary by distribution

### Permission Denied
```
Error: permission denied
```

**Solutions**:
- Tool installation may need sudo
- Run with appropriate privileges
- Check write permissions on home directory

### Dotfiles Script Fails
```
⚠️ Install script failed
```

**Solutions**:
- Check if script exists: `ls ~/src/dotfiles/fresh.sh`
- Run manually: `cd ~/src/dotfiles && bash ./fresh.sh`
- Check script has execution permissions
- May require specific shell (bash/zsh)

### SSH Key Registration Fails
```
❌ Failed to register SSH key
NOT AUTHENTICATED WITH GITHUB CLI
```

**Solutions**:
- Authenticate first: `gh auth login`
- Check network connectivity
- Verify GitHub account access

## Configuration

### Custom Dotfiles Repository
Currently hardcoded to: `https://github.com/aallbrig/dotfiles`

To use a different repository, modify the configuration in:
`pkg/computersetup/providers/<Platform>ConfigurationProvider.go`

### Custom Install Location
Dotfiles are cloned to: `~/src/dotfiles`

## Next Steps

After installation:
```bash
# Verify everything installed
allbctl bootstrap status

# Check your new environment
allbctl status

# Start coding!
cd ~/src
```
