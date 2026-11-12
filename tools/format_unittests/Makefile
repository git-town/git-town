RTA_VERSION = 0.43.3

RTA          = ../rta@$(RTA_VERSION)
ALPHASORT    = $(RTA) alphasort
DEADCODE     = $(RTA) deadcode
GOFUMPT      = $(RTA) gofumpt
GOLANGCILINT = $(RTA) golangci-lint

fix:  # runs all linters and auto-fixes
	$(ALPHASORT) fix
	$(GOFUMPT) -l -w .
	go run ../format_unittests/format_unittests.go
	go run ../format_self/format_self.go

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint:  # runs all linters
	$(DEADCODE) github.com/git-town/git-town/tools/format_unittests
	../ensure_no_files_with_dashes.sh
	$(GOLANGCILINT) cache clean
	$(GOLANGCILINT) run

test: unit lint  # runs all tests

unit:  # runs only the unit tests
	go test ./...

.DEFAULT_GOAL := help
