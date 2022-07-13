.PHONY: build
build:
	go build -v ./cmd/multiplexer

run: build
	./multiplexer

