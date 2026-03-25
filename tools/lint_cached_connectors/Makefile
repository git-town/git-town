.DEFAULT_GOAL := help

RTA_VERSION = 0.30.0  # run-that-app version to use
RTA = ../rta@$(RTA_VERSION)

fix:  # runs all linters and auto-fixes
	$(RTA) gofumpt -l -w .
	go run ../format_self/format_self.go

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint:  # runs all linters
	$(RTA) golangci-lint cache clean
	$(RTA) golangci-lint run
	$(RTA) deadcode github.com/git-town/git-town/tools/lint_cached_connectors

test: lint unit fix  # runs all tests

unit:  # runs the unit tests
	@env GOGC=off go test ./...
