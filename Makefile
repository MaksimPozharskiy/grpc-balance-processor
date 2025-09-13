.PHONY: generate build test run lint

generate:
	protoc -I api --go_out=. --go-grpc_out=. api/*.proto

build: generate
	go build -o bin/app ./cmd/app

test:
	go test -v ./...

run: build
	./bin/app

lint:
	golangci-lint run
