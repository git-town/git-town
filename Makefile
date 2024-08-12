RTA_VERSION = 0.6.1  # run-that-app version to use

# internal data and state
.DEFAULT_GOAL := help
RELEASE_VERSION := "15.1.0"
GO_BUILD_ARGS = LANG=C GOGC=off

build:  # builds for the current platform
	@go install -ldflags="-s -w"

cuke: build  # runs all end-to-end tests except the ones that mess up the output, best for development
	@env $(GO_BUILD_ARGS) skipmessyoutput=1 go test -v

cukeall: build  # runs all end-to-end tests
	@env $(GO_BUILD_ARGS) go test -v

cukethis: build  # runs the end-to-end tests that have a @this tag
	@env $(GO_BUILD_ARGS) cukethis=1 go test . -v -count=1

cukethiswin:  # runs the end-to-end tests that have a @this tag on Windows
	go install -ldflags "-X github.com/git-town/git-town/v15/internal/cmd.version=-dev -X github.com/git-town/git-town/v15/internal/cmd.buildDate=1/2/3"
	powershell -Command '$$env:cukethis=1 ; go test . -v -count=1'

cuke-prof: build  # creates a flamegraph for the end-to-end tests
	env $(GO_BUILD_ARGS) go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

cukewin: build  # runs all end-to-end tests on Windows
	go test . -v -count=1

dependencies: tools/rta@${RTA_VERSION}  # prints the dependencies between the internal Go packages as a tree
	@tools/rta depth . | grep git-town

docs: build tools/node_modules  # tests the documentation
	${CURDIR}/tools/node_modules/.bin/text-run --offline

fix: tools/rta@${RTA_VERSION} tools/node_modules  # runs all linters and auto-fixes
	go run tools/format_unittests/format_unittests.go
	go run tools/format_self/format_self.go
	tools/rta gofumpt -l -w .
	tools/rta dprint fmt
	tools/rta dprint fmt --config dprint-changelog.json
	${CURDIR}/tools/node_modules/.bin/prettier --write '**/*.yml'
	tools/rta shfmt -f . | grep -v tools/node_modules | grep -v '^vendor/' | xargs tools/rta shfmt --write
	tools/rta ghokin fmt replace features/

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

lint: tools/rta@${RTA_VERSION}  # lints the main codebase concurrently
	make --no-print-dir lint-smoke
	@tools/rta --available alphavet && go vet "-vettool=$(shell tools/rta --which alphavet)" $(shell go list ./... | grep -v internal/cmd)
	make --no-print-directory deadcode
	make --no-print-directory lint-structs-sorted
	git diff --check
	(cd tools/lint_steps && go build && ./lint_steps)
	${CURDIR}/tools/node_modules/.bin/gherkin-lint
	tools/rta actionlint
	tools/ensure_no_files_with_dashes.sh
	tools/rta shfmt -f . | grep -v 'tools/node_modules' | grep -v '^vendor/' | xargs tools/rta --optional shellcheck
	tools/rta golangci-lint cache clean
	tools/rta golangci-lint run

lint-all: lint tools/rta@${RTA_VERSION}  # runs all linters
	@echo lint tools/format_self
	@(cd tools/format_self && ../rta golangci-lint run)
	@echo lint tools/format_unittests
	@(cd tools/format_unittests && ../rta golangci-lint run)
	@echo lint tools/stats_release
	@(cd tools/stats_release && ../rta golangci-lint run)
	@echo lint tools/structs_sorted
	@(cd tools/structs_sorted && ../rta golangci-lint run)
	@echo lint tools/lint_steps
	@(cd tools/lint_steps && ../rta golangci-lint run)

lint-smoke: tools/rta@${RTA_VERSION}  # runs only the essential linters to get quick feedback after refactoring
	@tools/rta exhaustruct -test=false "-i=github.com/git-town/git-town.*" github.com/git-town/git-town/...
# @tools/rta ireturn --reject="github.com/git-town/git-town/v15/internal/gohacks/prelude.Option" github.com/git-town/git-town/...

lint-structs-sorted:
	@(cd tools/structs_sorted && go build) && ./tools/structs_sorted/structs_sorted

smoke: build  # run the smoke tests
	@env $(GO_BUILD_ARGS) smoke=1 go test . -v -count=1

smokewin: build  # runs the Windows smoke tests
	@env smoke=1 go test . -v -count=1

stats: tools/rta@${RTA_VERSION}  # shows code statistics
	@find . -type f | grep -v './tools/node_modules' | grep -v '\./vendor/' | grep -v '\./.git/' | grep -v './website/book' | xargs tools/rta scc

stats-release:  # displays statistics about the changes since the last release
	@(cd tools/stats_release && go build && ./stats_release v${RELEASE_VERSION})

test: fix docs unit lint-all cuke  # runs all the tests
.PHONY: test

test-go:  # smoke tests while working on the Go code
	@make --no-print-directory build &
	@make --no-print-directory unit &
	@make --no-print-directory deadcode &
	@make --no-print-directory lint

todo:  # displays all TODO items
	@git grep --color=always --line-number TODO ':!vendor' | grep -v Makefile

unit: build  # runs only the unit tests for changed code
	@env GOGC=off go test -timeout 30s ./internal/... ./pkg/... ./test/... ./tools/format_self/... ./tools/format_unittests/... ./tools/stats_release/... ./tools/structs_sorted/... ./tools/lint_steps/...

unit-all: build  # runs all the unit tests
	env GOGC=off go test -count=1 -timeout 60s ./internal/... ./pkg/... ./test/...

unit-race: build  # runs all the unit tests with race detector
	env GOGC=off go test -count=1 -timeout 60s -race ./internal/... ./pkg/... ./test/...

update: tools/rta@${RTA_VERSION}  # updates all dependencies
	go get -u ./...
	go mod tidy
	go work vendor
	(cd tools && yarn upgrade --latest)
	tools/rta --update

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

deadcode: tools/rta@${RTA_VERSION}
	@tput bold || true
	@tput setaf 1 || true
	@tools/rta deadcode github.com/git-town/git-town/tools/format_self &
	@tools/rta deadcode github.com/git-town/git-town/tools/format_unittests &
	@tools/rta deadcode github.com/git-town/git-town/tools/stats_release &
	@tools/rta deadcode github.com/git-town/git-town/tools/structs_sorted &
	@tools/rta deadcode github.com/git-town/git-town/tools/lint_steps &
	@tools/rta deadcode -test github.com/git-town/git-town/v15 | grep -v BranchExists \
	                                                           | grep -v 'Create$$' \
	                                                           | grep -v CreateFile \
	                                                           | grep -v CreateGitTown \
	                                                           | grep -v EmptyConfigSnapshot \
	                                                           | grep -v FileExists \
	                                                           | grep -v FileHasContent \
	                                                           | grep -v FilterErr \
	                                                           | grep -v IsGitRepo \
	                                                           | grep -v Memoized.AsFixture \
																														 | grep -v NewCommitMessages \
	                                                           | grep -v NewSHAs \
	                                                           | grep -v NewSet \
	                                                           | grep -v Paniced \
	                                                           | grep -v Set.Add \
	                                                           || true
	@tput sgr0 || true

tools/rta@${RTA_VERSION}:
	@rm -f tools/rta*
	@(cd tools && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh)
	@mv tools/rta tools/rta@${RTA_VERSION}
	@ln -s rta@${RTA_VERSION} tools/rta

tools/node_modules: tools/yarn.lock
	@echo "Installing Node based tools"
	@cd tools && yarn install
	@touch tools/node_modules  # update timestamp of the node_modules folder so that Make doesn't re-install it on every command
