---
weight: 6
title: "Git"
---

# Status Git

Display git global configuration including user name, email, and default editor.

## Usage

```bash
allbctl status git
```

## Output

```
Git Global Configuration:

  User Name:  Jane Developer
  User Email: jane@example.com
  Editor:     vim
```

If git is not configured:
```
Git Global Configuration:

  User Name:  (not set)
  User Email: (not set)
  Editor:     (not set)
```

## Configuration Shown

| Field | Source |
|-------|--------|
| User Name | `git config --global user.name` |
| User Email | `git config --global user.email` |
| Editor | `git config --global core.editor` |

## Notes

- Only reads global git configuration (`~/.gitconfig`), not per-repo settings
- Shows "(not set)" for unconfigured values
- Does not show credential helpers or other git settings

## Integration

Git configuration is shown inline in the `allbctl status` output. Run this subcommand to see it in isolation.
