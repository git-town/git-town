.DEFAULT_GOAL := help

build:  # builds the linter
	GOWORK=off go build -o lint_cached_connectors .

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint:  # runs all linters
	../rta golangci-lint run

test: build  # runs the linter
	./lint_cached_connectors

clean:  # removes the built binary
	rm -f lint_cached_connectors
