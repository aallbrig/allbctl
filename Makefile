LDFLAGS=-ldflags="-X 'github.com/aallbrig/allbctl/cmd.Version=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev")' -X 'github.com/aallbrig/allbctl/cmd.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")'"

install-dependencies:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/mitchellh/gox@latest
	go get -u github.com/fatih/color
	go get -u github.com/pkg/errors
	go mod download
	go mod vendor

install:
	go install $(LDFLAGS)

build:
	go build $(LDFLAGS) -o bin/allbctl main.go

install-local: build
	mkdir -p $(HOME)/go/bin
	cp bin/allbctl $(HOME)/go/bin/allbctl

# Docker targets
# ---------------------------------------------------------------------------
# docker-build: build the allbctl container image with version info baked in.
docker-build:
	docker build \
		--build-arg VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev") \
		--build-arg COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
		-t allbctl:latest \
		-t allbctl:$$(git describe --tags --always --dirty 2>/dev/null || echo "dev") \
		.

# docker-run: run allbctl with host namespaces and read-only mounts so that
# `allbctl status` inspects the real host (not the container).
# Pass ARGS to forward arguments, e.g.: make docker-run ARGS="status projects"
ARGS ?= status
docker-run: docker-build
	docker run --rm \
		--pid=host \
		--net=host \
		-v /proc:/proc:ro \
		-v /sys:/sys:ro \
		-v /etc:/etc:ro \
		-v $(HOME):$(HOME):ro \
		-e HOME=$(HOME) \
		allbctl:latest $(ARGS)

# Cross-platform builds — embed version/commit, no gox required
build-mac:
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o allbctl_darwin_amd64  main.go
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o allbctl_darwin_arm64  main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o allbctl_windows_amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o allbctl_windows_arm64.exe main.go

build-linux:
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o allbctl_linux_amd64   main.go
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o allbctl_linux_arm64   main.go

build-all: build-linux build-mac build-windows
	chmod +x allbctl_linux_* allbctl_darwin_*

# vagrant: boot Windows10 VM and run the allbctl smoke-test suite
# Requires: make build-windows first
vagrant-test-windows: build-windows
	vagrant up windows10
	vagrant provision windows10 --provision-with smoke-test

gen-docs:
	go run main.go gen-docs

check-docs:
	go run main.go gen-docs && git diff --exit-code hugo/site/content/docs/reference/

test:
	go test -v ./...

lint:
	golangci-lint run ./...

run:
	go run main.go

# Dev / CI workflow
# ---------------------------------------------------------------------------
# dev: start the "allbctl" tmux session (sites + watchers windows).
#      Requires tmuxinator and watchexec.
dev:
	tmuxinator start -p .tmuxinator/dev.yml

# ci: open the allbctl-ci tmux session that runs all quality gates in parallel
#     panes. Requires tmuxinator, golangci-lint, govulncheck, and gitleaks.
ci:
	tmuxinator start -p .tmuxinator/ci.yml

# ci-local: run all quality gates sequentially in the current shell (no tmux).
#           Suitable for scripted/headless CI environments.
ci-local:
	@echo "=== go fmt check ===" && \
	UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then echo "FAIL: unformatted files:" && echo "$$UNFORMATTED" && exit 1; fi && \
	echo "PASS: go fmt"
	@echo "=== go vet ===" && go vet ./... && echo "PASS: go vet"
	@echo "=== golangci-lint ===" && golangci-lint run ./... && echo "PASS: golangci-lint"
	@echo "=== go test ===" && go test -timeout 120s ./... && echo "PASS: go test"
	@echo "=== govulncheck ===" && govulncheck ./... && echo "PASS: govulncheck"
	@echo "=== check-docs ===" && make check-docs && echo "PASS: docs in sync"
	@echo "=== gitleaks ===" && gitleaks detect --source . --no-banner && echo "PASS: no secrets detected"
	@echo "All quality gates passed."