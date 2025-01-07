help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test:  # runs all tests
	go test ./...
	# prealloc suggests a premature optimisation.
	../rta golangci-lint run --disable prealloc
	../rta deadcode github.com/git-town/git-town/tools/tests_sorted
