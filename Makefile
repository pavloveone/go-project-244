APP_NAME=gendiff
BIN=bin/${APP_NAME}

.PHONY: build run lint test

build:
	mkdir -p bin
	go build -o ${BIN} ./cmd/${APP_NAME}

run:
	./${BIN} ${ARGS}

lint:
	golangci-lint run  ./...

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html