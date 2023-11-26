# dev tooling and versions
ALPHAVET_VERSION = 0.1.0
DEPTH_VERSION = 1.2.1
GOFUMPT_VERSION = 0.4.0
GOLANGCILINT_VERSION = 1.54.2
GORELEASER_VERSION = 1.22.1
RUN_THAT_APP_VERSION = 0.0.2
SCC_VERSION = 3.1.0
SHELLCHECK_VERSION = 0.9.0
SHFMT_VERSION = 3.6.0

# internal data and state
.DEFAULT_GOAL := help
TODAY = $(shell date +'%Y-%m-%d')
DEV_VERSION := $(shell git describe --tags 2> /dev/null || git rev-parse --short HEAD)
RELEASE_VERSION := "10.0.3"
GO_BUILD_ARGS = LANG=C GOGC=off

build:  # builds for the current platform
	go install -ldflags "-X github.com/git-town/git-town/v10/src/cmd.version=${DEV_VERSION}-dev -X github.com/git-town/git-town/v10/src/cmd.buildDate=${TODAY}"

cuke: build   # runs all end-to-end tests
	@env $(GO_BUILD_ARGS) go test . -v -count=1

cukethis: build   # runs the end-to-end tests that have a @this tag
	@env $(GO_BUILD_ARGS) cukethis=1 go test . -v -count=1

cuke-prof: build  # creates a flamegraph for the end-to-end tests
	env $(GO_BUILD_ARGS) go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

dependencies: tools/run-that-app@${RUN_THAT_APP_VERSION}  # prints the dependencies between the internal Go packages as a tree
	@tools/rta depth@${DEPTH_VERSION} . | grep git-town

docs: build tools/node_modules  # tests the documentation
	${CURDIR}/tools/node_modules/.bin/text-run --offline

fix: tools/alphavet_${ALPHAVET_VERSION} tools/run-that-app@${RUN_THAT_APP_VERSION} tools/node_modules  # auto-fixes lint issues in all languages
	git diff --check
	go run tools/format_unittests/format.go run
	go run tools/format_self/format.go run
	tools/rta gofumpt@${GOFUMPT_VERSION} -l -w .
	${CURDIR}/tools/node_modules/.bin/dprint fmt
	${CURDIR}/tools/node_modules/.bin/prettier --write '**/*.yml'
	tools/rta shfmt@${SHFMT_VERSION} -f . | grep -v tools/node_modules | grep -v '^vendor/' | xargs tools/run-that-app@${RUN_THAT_APP_VERSION} shfmt@${SHFMT_VERSION} --write
	tools/rta shfmt@${SHFMT_VERSION} -f . | grep -v tools/node_modules | grep -v '^vendor/' | xargs tools/run-that-app@${RUN_THAT_APP_VERSION} shellcheck@${SHELLCHECK_VERSION}
	${CURDIR}/tools/node_modules/.bin/gherkin-lint
	tools/rta golangci-lint@${GOLANGCILINT_VERSION} run
	tools/ensure_no_files_with_dashes.sh
	go vet "-vettool=tools/alphavet_${ALPHAVET_VERSION}" $(shell go list ./... | grep -v src/cmd | grep -v /v10/tools/)

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

stats: tools/run-that-app@${RUN_THAT_APP_VERSION}  # shows code statistics
	@find . -type f | grep -v './tools/node_modules' | grep -v '\./vendor/' | grep -v '\./.git/' | grep -v './website/book' | xargs tools/rta scc@${SCC_VERSION}

test: fix docs unit cuke  # runs all the tests
.PHONY: test

test-go: tools/alphavet_${ALPHAVET_VERSION} tools/run-that-app@${RUN_THAT_APP_VERSION}  # smoke tests for Go refactorings
	tools/rta gofumpt@${GOFUMPT_VERSION} -l -w . &
	make --no-print-directory build &
	tools/rta golangci-lint@${GOLANGCILINT_VERSION} run &
	go run tools/format_unittests/format.go test &
	go run tools/format_self/format.go test &
	@go vet "-vettool=tools/alphavet_${ALPHAVET_VERSION}" $(shell go list ./... | grep -v src/cmd | grep -v /v10/tools/) &
	make --no-print-directory unit

todo:  # displays all TODO items
	git grep --line-number TODO ':!vendor'

unit:  # runs only the unit tests for changed code
	@env GOGC=off go test -timeout 30s ./src/... ./test/...
	@go run tools/format_unittests/format.go test
	@go run tools/format_self/format.go test

unit-all:  # runs all the unit tests
	env GOGC=off go test -count=1 -timeout 60s ./src/... ./test/...

unit-race:  # runs all the unit tests with race detector
	env GOGC=off go test -count=1 -timeout 60s -race ./src/... ./test/...

update:  # updates all dependencies
	go get -u ./...
	go mod tidy
	go mod vendor
	(cd tools && yarn upgrade --latest)
	echo
	echo Please update the third-party tooling in the Makefile manually.

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

tools/alphavet_${ALPHAVET_VERSION}:
	@echo "Installing alphavet ${ALPHAVET_VERSION} ..."
	@env GOBIN="$(CURDIR)/tools" go install github.com/skx/alphavet/cmd/alphavet@latest
	@mv tools/alphavet tools/alphavet_${ALPHAVET_VERSION}

tools/run-that-app@${RUN_THAT_APP_VERSION}:
	@echo "Installing run-that-app ${RUN_THAT_APP_VERSION} ..."
	@rm -f tools/run-that-app* tools/rta
	@(cd tools && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh)
	@mv tools/run-that-app tools/run-that-app@${RUN_THAT_APP_VERSION}
	@ln -s run-that-app@${RUN_THAT_APP_VERSION} tools/rta

tools/node_modules: tools/yarn.lock
	@echo "Installing Node based tools"
	@cd tools && yarn install
	@touch tools/node_modules  # update timestamp of the node_modules folder so that Make doesn't re-install it on every command
