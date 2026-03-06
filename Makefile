install-dependencies:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/mitchellh/gox@latest
	go get -u github.com/fatih/color
	go get -u github.com/pkg/errors
	go mod download
	go mod vendor

install:
	go install -ldflags="-X 'github.com/aallbrig/allbctl/cmd.Version=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev")' -X 'github.com/aallbrig/allbctl/cmd.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")'"

build:
	go build -ldflags="-X 'github.com/aallbrig/allbctl/cmd.Version=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev")' -X 'github.com/aallbrig/allbctl/cmd.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")'" -o bin/allbctl main.go

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

build-mac:
	gox -osarch="darwin/amd64" -osarch="darwin/arm64"

build-windows:
	gox -osarch="windows/amd64" -osarch="windows/arm64"

build-linux:
	gox -osarch="linux/amd64" -osarch="linux/arm64"

build-all:
	gox -osarch="linux/amd64" -osarch="linux/arm64" -osarch="windows/amd64" -osarch="windows/arm64" -osarch="darwin/amd64" -osarch="darwin/arm64"
	chmod +x allbctl_*

test:
	go test -v ./...

lint:
	golangci-lint run ./...

run:
	go run main.go