SHELL = bash

default: test

.PHONY: test
test:
	@echo "--> Running Tests ..."
	@go test -v -race ./...

vet:
	@echo "--> Vet Go sources ..."
	@go vet ./...

generate:
	@echo "--> Go generate ..."
	@go generate ./...

changes: generate
	@echo "--> Checking for source diffs ..."
	@go mod tidy
	@go fmt ./...
	@./scripts/changes.sh
