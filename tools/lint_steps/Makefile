RTA_VERSION = 0.6.1  # run-that-app version to use

fix:  # runs all linters and auto-fixes
	../rta gofumpt -l -w .
	go run ../format_unittests/format_unittests.go
	go run ../format_self/format_self.go


help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint:  # runs all linters
	../rta golangci-lint run
	../rta deadcode github.com/git-town/git-town/tools/structs_sorted
	@../rta --available alphavet && go vet "-vettool=$(shell ../rta --which alphavet)" $(shell go list ./...)
	(cd ../structs_sorted && go build) && ../structs_sorted/structs_sorted
	../ensure_no_files_with_dashes.sh
	../rta golangci-lint cache clean
	../rta golangci-lint run

test: unit lint  # runs all tests

unit:  # runs only the unit tests
	go test ./...
