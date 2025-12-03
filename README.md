## allbctl

`allbctl` is short for allbrightctl and represents a command line interface for computer operations that I (Andrew Allbright) do. This is meant to be a CLI that is used by myself.

### Docs
```bash
# Help
allbctl --help
allbctl new-unity-project --help

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

## My favorites
project_name=$(basename "`pwd`")

allbctl new-unity-project \
  --project-name "${project_name}" \
  --create-repository-directory false \
  --install-webgl-fullscreen-template
  

```

### Bootstrapping a New Machine

On a fresh Linux machine, you can bootstrap your dev environment:

```bash
# 1. Install the latest allbctl from GitHub
go install github.com/aallbrig/allbctl@latest

# 2. Check what needs to be set up
allbctl cs status

# 3. Run the setup
allbctl cs install

# The tool will:
# - Create ~/src directory if needed
# - Verify git is installed (prompts you to install if missing)
# - Clone your dotfiles repo
# - Run the dotfiles install script
```

### Build
```bash
make install-dependencies
make build
```

### Tests
```bash
make lint
make test
```

### Install
```bash
make install
go install
```

### Contributing
Please reference the `CONTRIBUTING.md` file.

