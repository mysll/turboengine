.PHONY: all prepare build run gotool install clean help

BINARY_NAME={{.Name}}_service
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

all: proto gotool build

build:
	${ENV} go build -ldflags="-w -s" -o  ${BINARY_NAME}${BINARY_EXT} ${GO_FILE}

proto:
	@go test ./api/proto

entity:
	@go test ./internal/def

run:
	@go run ./

gotool:
	go fmt ./
	go vet ./

install:
	make proto build
	${MV} ${BINARY_NAME}${BINARY_EXT} ${BIN_DIR}

clean: 
	@if [ -f ${BINARY_NAME}${BINARY_EXT} ] ; then rm ${BINARY_NAME}${BINARY_EXT} ; fi
