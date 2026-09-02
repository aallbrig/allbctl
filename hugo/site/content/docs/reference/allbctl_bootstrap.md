---
title: "allbctl_bootstrap"
description: "auto-generated reference"
---

## allbctl bootstrap

Manage workstation bootstrap configuration

### Synopsis

Manage workstation bootstrap configuration including directories, tools, SSH keys, and dotfiles.

Available subcommands:
  status - Check current bootstrap status
  install - Apply bootstrap configuration to setup this machine
  reset - Reset bootstrap configuration

```
allbctl bootstrap [flags]
```

### Options

```
  -h, --help   help for bootstrap
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.allbctl.yaml)
      --debug           Enable debug telemetry (structured logs, traces, and metrics to stderr)
```

### SEE ALSO

* [allbctl](allbctl.md)	 - allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks
* [allbctl bootstrap install](allbctl_bootstrap_install.md)	 - Install and configure workstation bootstrap
* [allbctl bootstrap reset](allbctl_bootstrap_reset.md)	 - Reset workstation bootstrap configuration
* [allbctl bootstrap status](allbctl_bootstrap_status.md)	 - Check workstation bootstrap status
