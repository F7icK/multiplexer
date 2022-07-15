.PHONY: build
build:
	go build -v ./cmd/multiplexer

run: lint build
	./multiplexer

lint:
	golangci-lint run -c ./golangci.yml ./...