.DEFAULT_GOAL := help

fix:  # runs all linters and auto-fixes
	../rta gofumpt -l -w .
	go run ../format_self/format_self.go

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint:  # runs all linters
	../rta golangci-lint run
	../rta deadcode github.com/git-town/git-town/tools/lint_cached_connectors

test: lint  # runs the linter
