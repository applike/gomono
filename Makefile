
.PHONY: build test
build:
	go install

test:
	go test ./...
