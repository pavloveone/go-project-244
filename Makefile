APP_NAME=gendiff
BIN=bin/${APP_NAME}

.PHONY: build run

build:
	mkdir -p bin
	go build -o ${BIN} ./cmd/${APP_NAME}

run:
	./${BIN} ${ARGS}