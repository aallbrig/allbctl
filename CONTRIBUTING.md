### Contributing

Thank you so much for wanting to contribute! Please fork, create a change set (ideally with tests!), and submit a pull request.

## Local Development

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

### Pre-Commit Checks

To avoid CI/CD failures, **always run `make lint` before committing**. You can optionally set up a Git pre-commit hook to do this automatically:

```bash
# Create a pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
echo "Running linter..."
make lint
if [ $? -ne 0 ]; then
    echo "❌ Linting failed. Please fix the issues before committing."
    exit 1
fi
echo "✅ Linting passed!"
EOF

chmod +x .git/hooks/pre-commit
```

### Install Locally
```bash
make install
# OR
go install
```

