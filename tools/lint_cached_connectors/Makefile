RTA_VERSION = 0.43.3

RTA          = ../rta@$(RTA_VERSION)
DEADCODE     = $(RTA) deadcode
GOFUMPT      = $(RTA) gofumpt
GOLANGCILINT = $(RTA) golangci-lint

fix:  # runs all linters and auto-fixes
	$(GOFUMPT) -l -w .
	go run ../format_self/format_self.go

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint:  # runs all linters
	$(DEADCODE) github.com/git-town/git-town/tools/lint_cached_connectors
	$(GOLANGCILINT) cache clean
	$(GOLANGCILINT) run

test: lint unit fix  # runs all tests

unit:  # runs the unit tests
	@env GOGC=off go test ./...

.DEFAULT_GOAL := help
