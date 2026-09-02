---
title: "allbctl"
description: "auto-generated reference"
---

## allbctl

allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks

### Synopsis

allbctl (aka allbrightctl) is a CLI for Andrew Allbright specific tasks.

Example commands for allbctl:

$ allbctl bootstrap status
$ allbctl bootstrap install
$ allbctl status
$ allbctl status runtimes              # Show detected programming runtimes
$ allbctl status projects              # Show git repositories in ~/src
$ allbctl status list-packages         # Show package counts from all package managers
$ allbctl status db                    # Show detected databases and their status
$ allbctl status network               # Show network interface information
$ allbctl status containers            # Show container/virtualization info
$ allbctl status security              # Show SSH keys, GPG keys, and keyring
$ allbctl status systemctl             # Show systemd service counts
$ allbctl status git                   # Show git global configuration
$ allbctl status ports                 # Show listening ports
$ allbctl status cloud-native          # Show cloud CLI versions and profiles
$ allbctl status cloud-native aws      # Show detailed AWS resources by region
$ allbctl update                       # Update all detected package managers
$ allbctl update --dry-run             # Preview updates without executing
$ allbctl update --managers apt,npm    # Only update apt and npm


```
allbctl [flags]
```

### Options

```
      --config string   config file (default is $HOME/.allbctl.yaml)
      --debug           Enable debug telemetry (structured logs, traces, and metrics to stderr)
  -h, --help            help for allbctl
```

### SEE ALSO

* [allbctl bootstrap](allbctl_bootstrap.md)	 - Manage workstation bootstrap configuration
* [allbctl completion](allbctl_completion.md)	 - Generate the autocompletion script for the specified shell
* [allbctl status](allbctl_status.md)	 - Display system information (like neofetch)
* [allbctl update](allbctl_update.md)	 - Update and upgrade all detected package managers
* [allbctl version](allbctl_version.md)	 - Print the version number of allbctl
