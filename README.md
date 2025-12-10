## allbctl

`allbctl` is short for allbrightctl and represents a command line interface for computer operations that I (Andrew Allbright) do. This is meant to be a CLI that is used by myself.

### Docs
```bash
# Help
allbctl --help

# System status (like neofetch, but better!)
allbctl status                     # Shows comprehensive system information:
                                   # - Host info (OS, hostname, shell, terminal, CPU, GPU, memory)
                                   # - Network info (interfaces, router, connection type)
                                   # - Computer setup status (dotfiles, directories, tools)
                                   # - Package managers (available system & runtime package managers)
                                   # - Package counts per package manager

# List installed packages
allbctl list-packages              # Summary: just show counts per package manager (default)
allbctl list-packages --detail     # Full listing of all packages
allbctl list-packages -d           # Short version of --detail

# Computer setup (bootstrap development environment)
allbctl computer-setup status      # Check what's set up and what's missing
allbctl computer-setup install     # Install/configure dev environment
allbctl cs status                  # Short alias for status
allbctl cs install                 # Short alias for install

# What computer-setup does:
# ✅ Ensures ~/src directory exists
# ✅ Verifies git is installed
# ✅ Clones dotfiles from https://github.com/aallbrig/dotfiles
# ✅ Runs dotfiles install script (./fresh.sh) which sets up:
#    - oh-my-zsh
#    - Symlinks for .zshrc, .gitconfig, .vimrc, .tmux.conf, .ssh/config
#    - zsh as default shell (on Linux)
# ✅ Checks tools referenced in shell config files (.zshrc, .bashrc, etc.)
#    - Detects commands in $(command), source <(command), which/command -v, eval patterns
#    - Shows which tools are INSTALLED (green) vs MISSING (red)
#    - Groups by config file with $HOME path notation
#    - Works across Linux, macOS, and Windows

```

### Features

#### System Status
The `status` command provides a comprehensive view of your system, similar to neofetch but tailored for development:
- **Host Information**: OS, hostname, shell, terminal, CPU, GPU(s), memory, hardware details
- **Network Information**: Network interfaces with IPs, router IP, connection type (WiFi/Ethernet)
- **Computer Setup Status**: Dotfiles location, required directories, installed tools
- **Package Managers**: Detects available package managers on your system
  - **System**: apt, dnf, yum, pacman, snap, flatpak, zypper, apk, homebrew, macports, chocolatey, winget, scoop, nix
  - **Runtime**: npm, pip, gem, cargo, composer, maven, gradle
  - **WSL** (Windows only): Detects WSL availability and package managers inside WSL
- **Package Counts**: Summary of installed packages per detected package manager

#### Package Management
Multi-platform package detection supporting:
- **Linux**: apt, dnf, yum, pacman, snap, flatpak, zypper, apk, nix, homebrew
- **macOS**: homebrew, macports, nix
- **Windows**: chocolatey, winget, scoop, plus WSL package managers
- **Runtime**: npm, pip, gem, cargo, composer, maven, gradle (all platforms)

#### Computer Setup Automation
Automate bootstrapping a new development machine with dotfiles and configurations.

**Features:**
- Creates required directories (`~/src`, `~/bin`)
- Verifies essential tools are installed (git, etc.)
- Clones and runs your dotfiles setup
- **Shell Config Tool Detection**: Automatically scans your shell configuration files (.zshrc, .bashrc, .bash_profile, .profile) to find tool dependencies
  - Extracts commands from: `$(command)`, `` `command` ``, `source <(command)`, `which command`, `command -v`, `eval "$(command)"`
  - Reports which tools are INSTALLED (green) vs MISSING (red)
  - Groups output by config file with `$HOME` paths for portability
  - Filters out common shell builtins to focus on external tools
  - OS-agnostic: works on Linux, macOS, and Windows

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

# 4. Run the setup
allbctl cs install

# The tool will:
# - Create ~/src directory if needed
# - Verify git is installed (prompts you to install if missing)
# - Clone your dotfiles repo
# - Run the dotfiles install script
```

### Contributing
Please reference the `CONTRIBUTING.md` file.

