install-dependencies:
	go get -u golang.org/x/lint/golint
	go get -u github.com/mitchellh/gox
	go get -u github.com/fatih/color
	go get -u github.com/pkg/errors
	go mod download
	go mod vendor

install:
	go install

build:
	go build -o bin/allbctl main.go

build-docker:
	docker build --tag aallbrig/allbctl .

build-mac:
	gox -osarch="darwin/amd64"

build-windows:
	gox -osarch="windows/amd64"

build-linux:
	gox -osarch="linux/amd64"

build-all:
	gox -osarch="linux/amd64" -osarch="windows/amd64" -osarch="darwin/amd64"
	chmod +x allbctl_*

test:
	go test -v ./...

lint:
	golint `go list ./... | grep -v vendor/`

run:
	go run main.go