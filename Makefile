install-dependencies:
	go get -u golang.org/x/lint/golint
	go get -u github.com/mitchellh/gox
	go get -u github.com/fatih/color
	go mod download

build:
	go build -o bin/allbctl main.go

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
	golint ./...

run:
	go run main.go