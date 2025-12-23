install-dependencies:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/mitchellh/gox@latest
	go get -u github.com/fatih/color
	go get -u github.com/pkg/errors
	go mod download
	go mod vendor

install:
	go install

build:
	go build -ldflags="-X 'github.com/aallbrig/allbctl/cmd.Version=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev")' -X 'github.com/aallbrig/allbctl/cmd.Commit=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")'" -o bin/allbctl main.go

install-local: build
	mkdir -p $(HOME)/go/bin
	cp bin/allbctl $(HOME)/go/bin/allbctl

build-docker:
	docker build --tag aallbrig/allbctl .

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