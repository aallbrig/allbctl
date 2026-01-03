---
weight: 3
title: "Packages"
---

# Status List-Packages

Show package counts from all detected package managers on the system.

## Usage

```bash
# Summary mode (default)
allbctl status list-packages

# Detailed mode (show all packages)
allbctl status list-packages --detail
allbctl status list-packages -d

# Show packages from specific manager
allbctl status list-packages apt
allbctl status list-packages npm
allbctl status list-packages flatpak
```

## Output

### Summary Mode
Shows package counts from all detected package managers:

```
dpkg:           2109 packages
apt:            1968 packages
flatpak:        2 packages
npm:            6 packages
pip:            107 packages
pipx:           1 packages
go:             0 packages
ollama:         3 models
vagrant:        1 VMs
vboxmanage:     1 VMs

Use --detail flag to see the full list of all installed packages.
Or specify a package manager: allbctl status list-packages <manager>
```

### Detailed Mode
Shows full list of all installed packages from all managers:

```bash
allbctl status list-packages --detail
```

Outputs complete package lists from each detected manager.

### Specific Manager
Shows packages from one package manager with the command to reproduce:

```bash
allbctl status list-packages apt
```

```
Packages installed via apt:
package1
package2
...

Command: apt list --installed
```

## Supported Package Managers

### System Package Managers
- **dpkg** - Debian package database (all .deb packages)
- **apt** - Advanced Package Tool (explicitly installed)
- **rpm** - Red Hat Package Manager
- **dnf** - Dandified YUM (Fedora)
- **yum** - Yellowdog Updater Modified (RHEL/CentOS)
- **pacman** - Arch Linux package manager
- **snap** - Ubuntu snap packages
- **flatpak** - Flatpak packages
- **brew** - Homebrew (macOS/Linux)
- **choco** - Chocolatey (Windows)
- **winget** - Windows Package Manager
- **scoop** - Scoop (Windows)

### Programming Package Managers
- **npm** - Node.js packages (global only)
- **pip** - Python packages (global only)
- **pipx** - Python applications
- **gem** - Ruby gems (global only)
- **cargo** - Rust crates (global only)
- **go** - Go modules (global only)

### Special Managers
- **ollama** - Ollama AI models
- **vagrant** - Vagrant virtual machines
- **vboxmanage** - VirtualBox virtual machines

## Package Counting

- **System managers** (apt, dnf, etc.): Only explicitly installed packages, not dependencies
- **Programming managers** (npm, pip, etc.): Only globally installed packages
- **Special managers**: Ollama shows models, Vagrant/VirtualBox show VMs

## Integration

The package counts are shown in the "Packages:" section of the main `allbctl status` command.
