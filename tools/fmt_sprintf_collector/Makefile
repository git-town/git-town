.DEFAULT_GOAL := help

fix:  # runs all linters and auto-fixes
	../rta gofumpt -l -w .
	go run ../format_self/format_self.go

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint:  # runs all linters
	../rta golangci-lint run
	../rta deadcode github.com/git-town/git-town/tools/fmt_sprintf_collector
	@../rta --available alphavet && go vet "-vettool=$(shell ../rta --which alphavet)" $(shell go list ./...)
	../ensure_no_files_with_dashes.sh
	../rta golangci-lint cache clean
	../rta golangci-lint run

test: lint  # runs all tests
