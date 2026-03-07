---
weight: 9
title: "Security"
---

# Status Security

Display information about SSH keys, GPG keys, and the kernel keyring.

## Usage

```bash
allbctl status security
```

## Output

```
Security/Authentication Status:

  SSH Keys (loaded in agent):
    - 4096 SHA256:ngrORQUuSsaBf... user@hostname (RSA)
    - 256 SHA256:l3jOmOEHWpGVr... /home/user/.ssh/id_ed25519 (ED25519)

  GPG Keys:
    - sec  4096R/ABC12345  2023-01-01  Jane Developer <jane@example.com>

  Kernel Keyring:
    3 keys in user keyring
```

No SSH keys loaded:
```
Security/Authentication Status:

  SSH Keys (loaded in agent):
    No SSH keys loaded in agent

  GPG Keys:
    No GPG keys found

  Kernel Keyring:
    0 keys in user keyring
```

## Information Displayed

### SSH Keys
- Lists SSH keys currently loaded in `ssh-agent`
- Shows key size (bits), fingerprint (truncated), comment, and algorithm
- Detected via `ssh-add -l`

### GPG Keys
- Lists GPG secret keys available in the user's keyring
- Shows key type, ID, creation date, and user ID
- Detected via `gpg --list-secret-keys`

### Kernel Keyring
- Count of keys in the user's kernel keyring
- Detected via `keyctl list @u` (Linux only)
- Shows key count only (not key contents)

## Notes

- SSH keys are only shown if `ssh-agent` is running and has keys loaded
- GPG requires `gpg` or `gpg2` to be installed
- Kernel keyring count is Linux-specific; shows 0 on macOS/Windows

## Integration

Security information is shown in the `allbctl status` output. Run this subcommand to view it in isolation.
