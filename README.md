## allbctl

`allbctl` is short for allbrightctl and represents a command line interface for computer operations that I (Andrew Allbright) do. This is meant to be a CLI that is used by myself.

### Docs
```bash
# Help
allbctl --help

# System status (like neofetch, but better!)
allbctl status                     # Shows comprehensive system information:
allbctl cs status                  # Same as above (alias)
                                   # - User@hostname header with separator
                                   # - OS, kernel, uptime, host/virtualization info
                                   # - Package counts (dpkg, rpm, flatpak, snap, brew, etc.)
                                   # - Shell, terminal, CPU, GPU(s), memory
                                   # - Computer setup status (dotfiles, directories, tools)

# List installed packages
allbctl list-packages              # Summary: just show counts per package manager (default)
allbctl list-packages --detail     # Full listing of all packages
allbctl list-packages -d           # Short version of --detail
allbctl list-packages apt          # List only apt packages (shows command for copy/paste)
allbctl list-packages npm          # List only npm packages (shows command for copy/paste)
allbctl list-packages flatpak      # List only flatpak packages (shows command for copy/paste)

# Detect runtimes (programming languages, databases, cloud tools)
allbctl runtimes                   # Shows detected development runtimes with versions:
                                   # - Languages: Node.js, Go, PHP, Java, Python
                                   # - Databases: MySQL, PostgreSQL, SQLite, MariaDB, MongoDB, Redis
                                   # - Cloud: Kubernetes, AWS CLI, Azure CLI, Google Cloud SDK
                                   # - HashiCorp: Terraform, Vault, Consul, Nomad

# Computer setup (bootstrap development environment)
allbctl computer-setup status      # Check what's set up and what's missing
allbctl computer-setup install     # Install/configure dev environment automatically
allbctl cs status                  # Short alias for status
allbctl cs install                 # Short alias for install

# What computer-setup does:
# ✅ Ensures ~/src directory exists
# ✅ **Automatically installs git** (cross-platform: apt/dnf/yum/pacman/zypper/apk on Linux, brew on macOS, winget/choco/scoop on Windows)
# ✅ **Automatically installs GitHub CLI (gh)** (cross-platform package manager support)
# ✅ **Generates SSH keys** (if missing) and **registers them with GitHub** automatically
# ✅ **Clones dotfiles** from https://github.com/aallbrig/dotfiles to ~/src/dotfiles
# ✅ Runs dotfiles install script (./fresh.sh) which sets up:
#    - oh-my-zsh
#    - Symlinks for .zshrc, .gitconfig, .vimrc, .tmux.conf, .ssh/config
#    - zsh as default shell (on Linux)
# ✅ Checks tools referenced in shell config files (.zshrc, .bashrc, etc.)
#    - Detects commands in $(command), source <(command), which/command -v, eval patterns
#    - Shows which tools are INSTALLED (green) vs MISSING (red)
#    - Groups by config file with $HOME path notation
# ✅ **Fully cross-platform** - Works on Linux, macOS, and Windows
# ✅ **Idempotent** - Safe to run multiple times, only installs what's missing

```

### Features

#### System Status
The `status` and `cs status` commands provide a neofetch-inspired view of your system:
- **User@Hostname Header**: Shows current user and machine name with separator line
- **OS Information**: Platform, version, architecture, kernel version
- **System Details**: Host/virtualization info, uptime
- **Package Counts**: Inline display of packages from detected package managers (dpkg, rpm, pacman, snap, flatpak, brew, choco, winget)
- **System Info**: Shell, terminal, CPU with core count and frequency, GPU(s), memory usage
- **Network**: Network interfaces, router IP, connection type
- **AI Agents**: Detected AI coding assistants with versions (e.g., "copilot (0.0.365), claude (2.0.76)")
- **Package Managers**: 
  - **System**: System package managers with versions (e.g., "apt (2.8.3), flatpak (1.14.6)")
  - **Language**: Language version managers with versions (e.g., "nvm (0.40.3), pyenv (2.3.0)")
  - **Runtime**: Runtime package managers with versions (e.g., "npm (11.6.2), pip (24.0)")
- **Runtimes**: Detected programming language runtimes with versions (e.g., "Python (3.12.3), Node.js (24.11.1)")
- **Computer Setup Status**: Dotfiles location, required directories, installed tools, SSH configuration

##### Supported AI Agents
The `status` command detects the following AI coding assistants:
- **GitHub Copilot CLI** (`copilot`)
- **Claude Code** (`claude`)
- **Cursor AI** (`cursor`)
- **Aider** (`aider`)
- **Continue.dev** (`continue`)
- **Cody** (Sourcegraph) (`cody`)
- **Tabby** (`tabby`)
- **Ollama** (`ollama`)
- **Amazon CodeWhisperer** (`codewhisperer`)

##### Supported Language Version Managers
The `status` command detects the following language version managers:
- **NVM** (Node Version Manager) (`nvm`)
- **pyenv** (Python) (`pyenv`)
- **rbenv** (Ruby) (`rbenv`)
- **jenv** (Java) (`jenv`)
- **rustup** (Rust) (`rustup`)
- **asdf** (Universal version manager) (`asdf`)
- **SDKMAN** (Software Development Kit Manager) (`sdkman`)

#### Package Management
Multi-platform package detection supporting:
- **Linux**: dpkg, rpm, apt, dnf, yum, pacman, snap, flatpak, zypper, apk, nix, homebrew
- **macOS**: homebrew, macports, nix
- **Windows**: chocolatey, winget, scoop, plus WSL package managers
- **Runtime**: npm, pip, pipx, gem, cargo, composer, maven, gradle, go (all platforms)

##### Supported Package Managers

| Package Manager | Linux | macOS | Windows | Type |
|----------------|-------|-------|---------|------|
| **apt** | ✅ | ❌ | ❌ | System |
| **dpkg** | ✅ | ❌ | ❌ | System |
| **rpm** | ✅ | ❌ | ❌ | System |
| **dnf** | ✅ | ❌ | ❌ | System |
| **yum** | ✅ | ❌ | ❌ | System |
| **pacman** | ✅ | ❌ | ❌ | System |
| **snap** | ✅ | ❌ | ❌ | System |
| **flatpak** | ✅ | ❌ | ❌ | System |
| **brew** | ✅ | ✅ | ❌ | System |
| **choco** | ❌ | ❌ | ✅ | System |
| **winget** | ❌ | ❌ | ✅ | System |
| **scoop** | ❌ | ❌ | ✅ | System |
| **npm** | ✅ | ✅ | ✅ | Runtime |
| **pip** | ✅ | ✅ | ✅ | Runtime |
| **pipx** | ✅ | ✅ | ✅ | Runtime |
| **gem** | ✅ | ✅ | ✅ | Runtime |
| **cargo** | ✅ | ✅ | ✅ | Runtime |
| **go** | ✅ | ✅ | ✅ | Runtime |

**Usage:**
- `allbctl list-packages` - Summary of all detected package managers
- `allbctl list-packages <manager>` - List packages for a specific manager (e.g., `apt`, `npm`, `flatpak`)
  - Displays the underlying command for easy copy/paste (e.g., "Command: apt-mark showmanual")
- `allbctl list-packages --detail` - Full listing of all packages from all managers

#### Computer Setup Automation
Fully automate bootstrapping a new development machine with dotfiles and essential tools.

**Features:**
- **Cross-platform tool installation**: Automatically installs git and GitHub CLI using native package managers
  - **Linux**: apt, dnf, yum, pacman, zypper, apk
  - **macOS**: Homebrew
  - **Windows**: winget, chocolatey, scoop
- **SSH key management**: Generates SSH keys if missing and automatically registers them with GitHub
- Creates required directories (`~/src`, `~/bin`)
- **Dotfiles integration**: Clones your dotfiles repo to `~/src/dotfiles` and runs install script
- **Shell Config Tool Detection**: Automatically scans your shell configuration files (.zshrc, .bashrc, .bash_profile, .profile) to find tool dependencies
  - Extracts commands from: `$(command)`, `` `command` ``, `source <(command)`, `which command`, `command -v`, `eval "$(command)"`
  - Reports which tools are INSTALLED (green) vs MISSING (red)
  - Groups output by config file with `$HOME` paths for portability
  - Filters out common shell builtins to focus on external tools
  - OS-agnostic: works on Linux, macOS, and Windows
- **Idempotent operations**: Safe to run multiple times, only installs what's missing

### Bootstrapping a New Machine

On a fresh machine (Linux, macOS, or Windows), you can bootstrap your dev environment:

```bash
# 1. Download the latest release from GitHub
# Visit https://github.com/aallbrig/allbctl/releases/latest
# Download the appropriate binary for your platform:
#   - allbctl-linux-amd64 for Linux
#   - allbctl-darwin-amd64 for macOS (Intel)
#   - allbctl-darwin-arm64 for macOS (Apple Silicon)
#   - allbctl-windows-amd64.exe for Windows
# Make it executable (Linux/macOS): chmod +x allbctl-*
# Move to your PATH (e.g., ~/bin or /usr/local/bin)

# 2. Check your system status
allbctl status

# 3. Check what needs to be set up
allbctl cs status

# 4. Run the automated setup
allbctl cs install

# The tool will automatically:
# - Install git if not present (using your system's package manager)
# - Install GitHub CLI if not present
# - Generate SSH keys if missing
# - Register your SSH keys with GitHub (requires 'gh auth login' first)
# - Create ~/src directory if needed
# - Clone your dotfiles repo to ~/src/dotfiles
# - Run the dotfiles install script (./fresh.sh)
# - All operations are idempotent - safe to run multiple times
```

### Additional Commands

#### Runtimes Detection
The `runtimes` command detects programming languages, databases, and cloud tools installed on your system:

```bash
allbctl runtimes
```

**Detects:**
- **Programming Languages**: Node.js, Go, PHP, Java, Python
- **SQL Databases**: MySQL, PostgreSQL, SQLite, MariaDB, SQL Server, Oracle
- **NoSQL Databases**: MongoDB, Redis, Cassandra
- **Kubernetes & Cloud**: kubectl, AWS CLI, Azure CLI, Google Cloud SDK
- **HashiCorp Tools**: Terraform, Vault, Consul, Nomad

#### Reset Configuration
The `reset` command resets your machine configuration:

```bash
allbctl reset
```

This will display system information and reset the computer setup configuration to its initial state.

### Testing

#### Windows VM Testing
To test allbctl on Windows 10 using Vagrant:

```bash
# Build Windows binary
make build-windows

# Start Windows 10 VM
vagrant up windows10

# The VM will boot with GUI - log in and open PowerShell
# Test the bootstrap sequence:
cd C:\allbctl-test
.\allbctl.exe bootstrap status    # Check initial state
.\allbctl.exe bootstrap install   # Install components
.\allbctl.exe bootstrap status    # Verify installation

# Cleanup
vagrant halt windows10      # Stop VM
vagrant destroy windows10   # Remove VM
```

See [test/windows-vm-test.md](test/windows-vm-test.md) for detailed testing instructions.

### Contributing
Please reference the `CONTRIBUTING.md` file.
