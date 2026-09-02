---
title: "allbctl_bootstrap_install"
description: "auto-generated reference"
---

## allbctl bootstrap install

Install and configure workstation bootstrap

### Synopsis

Install and configure workstation bootstrap including directories, tools, SSH keys, and dotfiles.

By default, SSH key generation and GitHub registration are SKIPPED.
Use --register-ssh-keys flag to enable SSH key generation and GitHub registration.

```
allbctl bootstrap install [flags]
```

### Options

```
  -h, --help                help for install
      --register-ssh-keys   Generate SSH keys and register with GitHub (requires gh CLI)
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.allbctl.yaml)
      --debug           Enable debug telemetry (structured logs, traces, and metrics to stderr)
```

### SEE ALSO

* [allbctl bootstrap](allbctl_bootstrap.md)	 - Manage workstation bootstrap configuration
