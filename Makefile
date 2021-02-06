install-dependencies:
	go mod download

build:
	go build -o bin/allbctl main.go

build-mac:
	gox -osarch="darwin/amd64"

build-windows:
	gox -osarch="windows/amd64"

build-linux:
	gox -osarch="linux/amd64"

test:
	go test -v ./...

lint:
	golint ./...

run:
	go run main.go