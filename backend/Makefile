.PHONY: proto build run clean test

PROTO_DIR := proto
GEN_DIR := gen

proto:
	protoc --go_out=. --go_opt=module=github.com/clarity/backend \
		--go-grpc_out=. --go-grpc_opt=module=github.com/clarity/backend \
		proto/*.proto

build:
	go build -o bin/clarity-backend main.go

run:
	go run main.go

test:
	go test ./... -v

clean:
	rm -rf bin/
	rm -rf gen/

deps:
	go mod download
	go mod tidy

all: proto deps build

.DEFAULT_GOAL := all
