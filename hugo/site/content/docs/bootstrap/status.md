---
weight: 1
title: "Bootstrap Status"
---

# Bootstrap Status

Check the status of your workstation bootstrap configuration.

## Usage

```bash
allbctl bootstrap status
```

## What It Shows

### Expected Directories
Shows whether required directories exist:
```
Expected Directories
-----
PRESENT /home/user/src
```

### Required Tools
Shows which development tools are installed:
```
Required Tools
-----
INSTALLED git
INSTALLED gh
```

Or if missing:
```
Required Tools
-----
INSTALLED git
NOT FOUND gh
```

### SSH Configuration
Shows SSH key status and GitHub registration:
```
SSH Configuration
-----
SSH KEY REGISTERED WITH GITHUB
```

Possible states:
- `SSH KEY REGISTERED WITH GITHUB` - Key exists and registered
- `SSH KEY NOT REGISTERED WITH GITHUB` - Key exists but not registered
- `SSH KEY NOT FOUND` - No SSH key exists
- `NOT AUTHENTICATED WITH GITHUB CLI` - gh CLI not authenticated

### Dotfiles
Shows whether dotfiles repository is cloned:
```
Dotfiles
-----
CLONED /home/user/src/dotfiles
```

Or if not:
```
Dotfiles
-----
NOT CLONED /home/user/src/dotfiles
```

### Shell Config Tools
Shows tools referenced in shell config files and their status:
```
Shell Config Tools
-----
$HOME/.zshrc:
INSTALLED kubectl
MISSING   aws_completer
INSTALLED tmux
```

## Example Output

Full example of bootstrap status:

```
Workstation Bootstrap Status:

  Expected Directories
  -----
  PRESENT /home/user/src

  Required Tools
  -----
  INSTALLED git
  INSTALLED gh

  SSH Configuration
  -----
  SSH KEY REGISTERED WITH GITHUB

  Dotfiles
  -----
  CLONED /home/user/src/dotfiles

  Shell Config Tools
  -----
  $HOME/.zshrc:
  INSTALLED kubectl
  INSTALLED tmux
  MISSING   aws_completer
```

## Use Cases

### Before Installing
Check what's missing before running install:
```bash
allbctl bootstrap status
```

### After Installing
Verify everything was installed correctly:
```bash
allbctl bootstrap install
allbctl bootstrap status
```

### Regular Checks
Periodically check your setup status:
```bash
allbctl bootstrap status
```

## Color Coding

- **Green (PRESENT/INSTALLED/REGISTERED)** - Component is set up correctly
- **Red (NOT FOUND/MISSING)** - Component is missing
- **Yellow** - Warning or partial state

## Next Steps

If status shows missing components:
```bash
# Install missing components
allbctl bootstrap install

# Or with SSH keys
allbctl bootstrap install --register-ssh-keys
```
