.PHONY: all test build clean

all: clean test build

build: 
	mkdir -p build
	go build -o build ./...

test:
	go test -v -coverprofile=tests/results/cover.out ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

verify:
	golangci-lint run -c .golangci.yaml --deadline=30m

clean:
	rm -rf build/*
	go clean ./...
