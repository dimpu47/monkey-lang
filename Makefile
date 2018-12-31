.PHONY: build test deps clean

all: build
	@./monkey

deps:
	@go get ./...

build:
	@go build -o monkey .

test:
	@go test -v -cover -coverprofile=coverage.txt -covermode=atomic ./...

clean:
	@git clean -f -d -X
