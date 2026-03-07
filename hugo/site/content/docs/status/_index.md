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

### Header
- Username@Hostname
- allbctl version and commit hash

### System Information
- Operating system and version
- Hostname
- Shell (e.g., `/usr/bin/zsh`, `C:\WINDOWS\system32\cmd.exe`)
- Terminal emulator with enhanced detection:
  - **Windows**: Windows Terminal, PowerShell Core, PowerShell, Command Prompt (cmd.exe), Git Bash, ConEmu
  - **Unix-like**: tmux, iTerm2, kitty, alacritty, konsole, gnome-terminal, and more
- CPU details (model, cores, architecture)
- GPU information
- Memory usage
- Hardware details

### Development Tools
- **Runtimes**: Detected programming languages (Python, Go, Node.js, Java, etc.)
- **Databases**: Detected database systems and their status
- **Package Managers**: System and programming package managers
- **Packages**: Package counts from all detected managers
  - **Performance**: Parallelized package detection streams results as they're counted
  - Package counts and update checks run concurrently for faster output
- **Projects**: Git repositories in ~/src directory
- **Cloud Native**: Cloud CLI tools (AWS CLI, gcloud, Azure CLI, kubectl) with profile counts

### Network Information
- Network interfaces
- IP addresses
- Router information
- Connection type

### Additional Information
- Browsers: Detected web browsers with versions
- AI Agents: Detected AI coding assistants

## Example Output

```
aallbright@unicorn-tp
allbctl v0.0.30 (commit abc1234)

OS:        linuxmint 22.2
Hostname:  unicorn-tp
Shell:     /usr/bin/zsh
Terminal:  tmux
CPU:
  Model:     13th Gen Intel(R) Core(TM) i7-1360P
  Arch:      x86_64
  Clock:     5.00 GHz
  Cores:     12 physical, 16 logical (2 threads/core)
GPU(s):
  Name:      NVIDIA RTX A500 Laptop GPU
  Memory:    4096 MiB
Memory:    33.3 GiB
Disks:     2 total (1007.0 GB)
Hardware:  unicorn-tp linuxmint
Runtimes:  Python (3.12.3 → 3.14.3), Java (21.0.10), Go (1.26.1), Node.js (24.11.1)
Databases: sqlite3 (3.45.1)

Network:
  Primary Interface: wlp0s20f3 (192.168.50.90/24)
    WiFi: MY_SSID @ 5.24 GHz (802.11ac/ax)
    Gateway: 192.168.50.1
  Connectivity:
    Internet: ✓ Connected

Ports:     24 listening (TCP: 18, UDP: 6)

Browsers:
  Chromium (145.0.7632.45), Firefox (147.0.3)

AI Agents:
  copilot (GitHub Copilot CLI 1.0.2.), claude (2.0.76)

Package Managers:
  System:   apt (2.8.3), flatpak (1.14.6)
  Runtime:  npm (11.6.2), pip (24.0), go (1.26.1)

Packages:
  dpkg:     2166 packages
  apt:      1991 packages (90 want updates)
  npm:      10 packages (5 want updates)
  pip:      110 packages

Cloud Native:
  kubectl (1.34.2) - 0 contexts ✗
  aws (2.33.1) - 1 profile ✓

Projects: 12 total (7 dirty)
  Last 5 recently touched:
    ~/src/allbctl*              aallbrig/allbctl        2026-03-07 16:19 EST
    ~/src/stock-market-words    aallbrig/stock-market-words  2026-03-07 11:40 EST
    ~/src/alpha_rush*           aallbrig/alpha_rush     2026-03-07 15:41 EST
    ~/src/assgen                aallbrig/assgen         2026-03-07 14:49 EST
    ~/src/treemand*             aallbrig/treemand       2026-03-07 14:30 EST
```



All status subcommands show specific sections of the main status output:

- [Runtimes](runtimes)
- [Projects](projects)
- [Packages](packages)
- [Databases](databases)
- [Cloud Native](cloud-native)
