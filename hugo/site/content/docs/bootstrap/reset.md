---
weight: 3
title: "Bootstrap Reset"
---

# Bootstrap Reset

Reset your bootstrap configuration by removing installed components.

## Usage

```bash
allbctl bootstrap reset
```

## What It Does

Removes bootstrap components in **reverse order** of installation:

1. **Shell Config Tools** - (tracking only, not removed)
2. **Dotfiles** - (manual removal recommended)
3. **SSH Configuration** - (manual removal recommended)
4. **Required Tools** - (manual removal recommended)
5. **Expected Directories** - (manual removal recommended)

## Important Notes

### Manual Removal Recommended

For safety, most components show instructions for manual removal rather than auto-removing:

```
❌ Cannot auto-uninstall SSH key registration - please remove manually from GitHub
❌ Cannot auto-uninstall dotfiles - please remove manually
```

This prevents accidental data loss.

### What Gets Auto-Removed

Currently, reset is primarily for testing and shows status but requires manual cleanup.

## Example Output

```bash
$ allbctl bootstrap reset
```

```
System Info
-----
OS: ubuntu 22.04
Hostname: workstation

Uninstalling: Shell Config Tools
❌ Cannot auto-uninstall shell config tracking

Uninstalling: Dotfiles
❌ Cannot auto-uninstall dotfiles from ~/src/dotfiles - please remove manually

Uninstalling: SSH Configuration
❌ Cannot auto-uninstall SSH key registration - please remove manually from GitHub

Uninstalling: Required Tools
❌ Cannot auto-uninstall gh - please remove manually
❌ Cannot auto-uninstall git - please remove manually

Uninstalling: Expected Directories
❌ Cannot auto-uninstall ~/src - please remove manually
```

## Manual Cleanup

### Remove Dotfiles
```bash
rm -rf ~/src/dotfiles
```

### Remove SSH Keys
1. Remove from GitHub:
   - Visit https://github.com/settings/keys
   - Delete the key for this machine

2. Remove local key:
   ```bash
   rm ~/.ssh/id_rsa ~/.ssh/id_rsa.pub
   ```

### Remove Tools

**Linux (apt)**:
```bash
sudo apt remove gh git
```

**Linux (dnf)**:
```bash
sudo dnf remove gh git
```

**macOS (brew)**:
```bash
brew uninstall gh git
```

**Windows (winget)**:
```powershell
winget uninstall GitHub.cli
winget uninstall Git.Git
```

### Remove Directories
```bash
# Only if empty or you're sure
rm -rf ~/src
```

## Why Manual Removal?

Bootstrap reset doesn't auto-remove components because:

1. **Data Safety** - Prevents accidental deletion of code in ~/src
2. **SSH Keys** - Removing from GitHub requires API access
3. **System Tools** - git/gh might be used by other applications
4. **Dotfiles** - May contain important customizations

## Testing Workflow

For testing bootstrap in a clean environment:

```bash
# Install
allbctl bootstrap install

# Check status
allbctl bootstrap status

# Reset (shows what to remove)
allbctl bootstrap reset

# Manual cleanup
rm -rf ~/src/dotfiles
# ... other manual steps

# Verify clean state
allbctl bootstrap status
```

## Use Cases

### Testing Bootstrap
Test installation in a VM:
```bash
# Install
allbctl bootstrap install

# Test
allbctl bootstrap status

# Clean up for next test
allbctl bootstrap reset
# Then manually remove components
```

### Troubleshooting Installation
If installation goes wrong:
```bash
# Reset to see what was installed
allbctl bootstrap reset

# Manually clean up problem components
# Then try again
allbctl bootstrap install
```

### Moving to Clean Slate
Before switching dotfiles repositories or major changes:
```bash
allbctl bootstrap reset
# Manually remove old setup
# Configure new setup
allbctl bootstrap install
```

## Future Enhancement

A future version might include:
- `--force` flag for actual removal
- Selective component removal
- Backup before removal
- More automated cleanup

## Alternatives

Instead of reset, you can:

### Selectively Remove Components
```bash
# Just remove dotfiles
rm -rf ~/src/dotfiles

# Just uninstall gh
sudo apt remove gh  # or brew uninstall gh
```

### Fresh OS Install
For complete clean slate:
- Reinstall operating system
- Or use a new VM/container

### Check and Reinstall
```bash
# Check what's there
allbctl bootstrap status

# Reinstall over existing (idempotent)
allbctl bootstrap install
```

## Related Commands

- [Bootstrap Status](status) - Check current state
- [Bootstrap Install](install) - Install components
