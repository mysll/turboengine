.PHONY: all prepare build run gotool install clean help

BINARY_NAME=cli
BIN_DIR=../../bin/
LAN_FILE=.go
GO_FILE:=main.go

all: prepare gotool build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o  ${BINARY_NAME}  ${GO_FILE}

prepare:
	go test ./api/proto

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make prepare build
	mv ${BINARY_NAME} ${BIN_DIR}

clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi
