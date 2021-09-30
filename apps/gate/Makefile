.PHONY: all prepare build run gotool install clean help

BINARY_NAME=gate_service
ifeq ($(OS),Windows_NT)
BINARY_EXT=.exe
ENV=
MV=move
else
BINARY_EXT=
ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64
MV=mv
endif
BIN_DIR=../../bin/
LAN_FILE=.go
GO_FILE:=main.go

all: prepare gotool build

build:
	${ENV} go build -ldflags="-w -s" -o  ${BINARY_NAME}${BINARY_EXT} ${GO_FILE}

prepare:
	@go test ./api/proto

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make prepare build
	${MV} ${BINARY_NAME}${BINARY_EXT} ${BIN_DIR}

clean: 
	@if [ -f ${BINARY_NAME}${BINARY_EXT} ] ; then rm ${BINARY_NAME}${BINARY_EXT} ; fi
