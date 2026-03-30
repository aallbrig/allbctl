---
weight: 1
title: "Installation"
---

# Installation

## Homebrew (macOS and Linux)

The easiest way to install allbctl on macOS or Linux:

```bash
brew tap aallbrig/tap
brew install allbctl
```

Or in one step:

```bash
brew install aallbrig/tap/allbctl
```

To update:

```bash
brew upgrade allbctl
```

## Chocolatey (Windows)

The easiest way to install allbctl on Windows:

```powershell
choco install allbctl
```

To update:

```powershell
choco upgrade allbctl
```

To uninstall:

```powershell
choco uninstall allbctl
```

## APT (Debian / Ubuntu)

Install allbctl from the official APT repository:

```bash
# Add the GPG key
curl -fsSL https://aallbrig.github.io/apt-repo/KEY.gpg | sudo gpg --dearmor -o /usr/share/keyrings/allbctl.gpg

# Add the repository
echo "deb [signed-by=/usr/share/keyrings/allbctl.gpg] https://aallbrig.github.io/apt-repo stable main" | sudo tee /etc/apt/sources.list.d/allbctl.list

# Install
sudo apt update
sudo apt install allbctl
```

To update:

```bash
sudo apt update && sudo apt upgrade allbctl
```

To uninstall:

```bash
sudo apt remove allbctl
```

## Download Latest Release

The easiest way to install allbctl is to download a pre-built binary from the releases page.

### Linux

```bash
# Download latest release for Linux
curl -LO https://github.com/aallbrig/allbctl/releases/latest/download/allbctl-linux-amd64

# Make executable
chmod +x allbctl-linux-amd64

# Move to PATH
sudo mv allbctl-linux-amd64 /usr/local/bin/allbctl

# Verify installation
allbctl --version
```

### macOS

```bash
# Download latest release for macOS (Intel)
curl -LO https://github.com/aallbrig/allbctl/releases/latest/download/allbctl-darwin-amd64

# Or for Apple Silicon
curl -LO https://github.com/aallbrig/allbctl/releases/latest/download/allbctl-darwin-arm64

# Make executable
chmod +x allbctl-darwin-*

# Move to PATH
sudo mv allbctl-darwin-* /usr/local/bin/allbctl

# Verify installation
allbctl --version
```

### Windows

```powershell
# Download latest release for Windows
# Visit: https://github.com/aallbrig/allbctl/releases/latest
# Download: allbctl-windows-amd64.exe

# Or use PowerShell
Invoke-WebRequest -Uri "https://github.com/aallbrig/allbctl/releases/latest/download/allbctl-windows-amd64.exe" -OutFile "allbctl.exe"

# Add to PATH or move to desired location
Move-Item allbctl.exe C:\Windows\System32\

# Verify installation
allbctl --version
```

## Prerequisites for Building from Source

- Go 1.20+ (for building from source)
- Git (for cloning repository)

## Install from Source

```bash
# Clone the repository
git clone https://github.com/aallbrig/allbctl.git
cd allbctl

# Build the binary
make build

# Install to /usr/local/bin (optional)
sudo cp bin/allbctl /usr/local/bin/

# Or add to PATH
export PATH=$PATH:$(pwd)/bin
```

## Build for Different Platforms

```bash
# Linux (default)
make build

# Windows
GOOS=windows GOARCH=amd64 go build -o allbctl.exe .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o allbctl .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o allbctl .
```

## Verify Installation

```bash
allbctl --version
allbctl --help
```

## Update

### From Release
Simply download and replace the binary with the latest version from [releases](https://github.com/aallbrig/allbctl/releases/latest).

### From Source
```bash
cd allbctl
git pull
make build
```

## Uninstall

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/allbctl

# Or just remove the repository
rm -rf ~/path/to/allbctl
```

