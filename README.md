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

## My favorites
project_name=$(basename "`pwd`")

allbctl new-unity-project \
  --project-name "${project_name}" \
  --create-repository-directory false \
  --install-webgl-fullscreen-template
  

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

