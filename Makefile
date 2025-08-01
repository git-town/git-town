RTA_VERSION = 0.17.0  # run-that-app version to use

# internal data and state
.DEFAULT_GOAL := help
RELEASE_VERSION := "21.4.0"
GO_TEST_ARGS = LANG=C GOGC=off BROWSER=

cuke: install  # runs all end-to-end tests in a way that looks nice during development
	@env $(GO_TEST_ARGS) messyoutput=0 go test -v
	@env $(GO_TEST_ARGS) messyoutput=1 go test -v

cukeall: install  # runs all end-to-end tests
	@env $(GO_TEST_ARGS) go test -v

cukethis: install  # runs the end-to-end tests that have a @this tag
	@env $(GO_TEST_ARGS) cukethis=1 go test . -v -count=1

cukethiswin:  # runs the end-to-end tests that have a @this tag on Windows
	go install -ldflags "-X github.com/git-town/git-town/v21/internal/cmd.version=-dev -X github.com/git-town/git-town/v21/internal/cmd.buildDate=1/2/3"
	powershell -Command '$$env:cukethis=1 ; go test . -v -count=1'

cuke-prof: install  # creates a flamegraph for the end-to-end tests
	env $(GO_TEST_ARGS) go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

cukeverbose: install  # run all tests in "verbose.feature" files
	@env $(GO_TEST_ARGS) verbose=1 go test . -v -count=1

cukewin: install  # runs all end-to-end tests on Windows
	go test . -v -count=1

dependencies: tools/rta@${RTA_VERSION}  # prints the dependencies between the internal Go packages as a tree
	@tools/rta depth . | grep git-town

docs: install tools/node_modules  # tests the documentation
	@tools/rta node tools/node_modules/.bin/text-runner --offline

fix: tools/rta@${RTA_VERSION}  # runs all linters and auto-fixes
	go run tools/format_unittests/format_unittests.go
	go run tools/format_self/format_self.go
	tools/rta gofumpt -l -w .
	tools/rta dprint fmt
	tools/rta dprint fmt --config dprint-changelog.json
	tools/rta shfmt -f . | grep -v tools/node_modules | grep -v '^vendor/' | xargs tools/rta shfmt --write
	tools/rta ghokin fmt replace features/
	tools/generate_opcodes_all.sh

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install:  # builds for the current platform
	@go install -ldflags="-s -w"

lint: tools/node_modules tools/rta@${RTA_VERSION}  # lints the main codebase concurrently
	make --no-print-directory lint-smoke
	make --no-print-directory alphavet
	make --no-print-directory deadcode
	make --no-print-directory lint-messages-sorted
	make --no-print-directory lint-messy-output
	make --no-print-directory lint-optioncompare
	make --no-print-directory lint-print-config
	make --no-print-directory lint-structs-sorted
	make --no-print-directory lint-tests-sorted
	make --no-print-directory lint-use-equal
	git diff --check
	(cd tools/lint_steps && go build && ./lint_steps)
	tools/rta actionlint
	tools/rta --from-source staticcheck ./...
	tools/ensure_no_files_with_dashes.sh
	tools/rta shfmt -f . | grep -v 'tools/node_modules' | grep -v '^vendor/' | xargs tools/rta --optional shellcheck
	tools/rta golangci-lint cache clean
	tools/rta golangci-lint run
	tools/rta node tools/node_modules/.bin/gherkin-lint

lint-all: lint tools/rta@${RTA_VERSION}  # runs all linters
	(cd website && make test)
	tools/rta govulncheck ./...
	@echo lint tools/format_self
	@(cd tools/format_self && make test)
	@echo lint tools/format_unittests
	@(cd tools/format_unittests && make test)
	@echo lint tools/lint_steps
	@(cd tools/lint_steps && make test)
	@echo lint tools/messages_sorted
	@(cd tools/messages_sorted && make lint)
	@echo lint tools/messy_output
	@(cd tools/messy_output && make test)
	@echo lint tools/optioncompare
	@(cd tools/optioncompare && make test)
	@echo lint tools/print_config_exhaustive
	@(cd tools/print_config_exhaustive && make test)
	@echo lint tools/stats_release
	@(cd tools/stats_release && make test)
	@echo lint tools/structs_sorted
	@(cd tools/structs_sorted && make test)
	@echo lint tools/tests_sorted
	@(cd tools/tests_sorted && make test)
	@echo lint tools/tests_sorted
	@(cd tools/tests_sorted && make test)
	@echo lint tools/use_equal
	@(cd tools/use_equal && make test)

alphavet:
	@tools/rta --available alphavet && go vet "-vettool=$(shell tools/rta --which alphavet)" $(shell go list ./... | grep -v internal/cmd)

lint-messages-sorted:
	@(cd tools/messages_sorted && go build) && ./tools/messages_sorted/messages_sorted

lint-messy-output:
	@(cd tools/messy_output && go build) && ./tools/messy_output/messy_output

lint-print-config:
	@(cd tools/print_config_exhaustive && go build) && ./tools/print_config_exhaustive/print_config_exhaustive

lint-optioncompare:
	@(cd tools/optioncompare && go build) && ./tools/optioncompare/optioncompare github.com/git-town/git-town/v21/...

lint-smoke: tools/rta@${RTA_VERSION}  # runs only the essential linters to get quick feedback after refactoring
	@tools/rta exhaustruct -test=false "-i=github.com/git-town/git-town.*" github.com/git-town/git-town/...
# @tools/rta ireturn --reject="github.com/git-town/git-town/v21/pkg/prelude.Option" github.com/git-town/git-town/...

lint-structs-sorted:
	@(cd tools/structs_sorted && go build) && ./tools/structs_sorted/structs_sorted

lint-tests-sorted:
	@(cd tools/tests_sorted && go build) && ./tools/tests_sorted/tests_sorted

lint-use-equal:
	@(cd tools/use_equal && go build) && ./tools/use_equal/use_equal

smoke: install  # run the smoke tests
	@env $(GO_TEST_ARGS) smoke=1 go test . -v -count=1

smokewin: install  # runs the Windows smoke tests
	@env smoke=1 go test . -v -count=1

stats: tools/rta@${RTA_VERSION}  # shows code statistics
	@find . -type f \
		| grep -v './tools/node_modules' \
		| grep -v '\./vendor/' \
		| grep -v '\./.git/' \
		| grep -v './website/book' \
		| xargs tools/rta scc

stats-release:  # displays statistics about the changes since the last release
	@(cd tools/stats_release && go build && ./stats_release v${RELEASE_VERSION})

test: fix docs unit lint-all cuke  # runs all the tests
.PHONY: test

test-go:  # smoke tests while working on the Go code
	@make --no-print-directory install &
	@make --no-print-directory unit &
	@make --no-print-directory deadcode &
	@make --no-print-directory lint

todo:  # displays all TODO items
	@git grep --color=always --line-number TODO ':!vendor' \
		| grep -v Makefile \
		| grep -v ':= context.'

UNIT_TEST_DIRS = \
	./internal/... \
	./pkg/... \
	./tools/format_self/... \
	./tools/format_unittests/... \
	./tools/lint_steps/... \
	./tools/messages_sorted/... \
	./tools/messy_output/... \
	./tools/stats_release/... \
	./tools/structs_sorted/... \
	./tools/tests_sorted/...

unit: install  # runs only the unit tests for changed code
	@env GOGC=off go test -timeout=30s $(UNIT_TEST_DIRS)

unit-all: install  # runs all the unit tests
	env GOGC=off go test -count=1 -shuffle=on -timeout=60s $(UNIT_TEST_DIRS)

unit-race: install  # runs all the unit tests with race detector
	env GOGC=off go test -count=1 -timeout 60s -race $(UNIT_TEST_DIRS)
	cd website && make --no-print-directory unit

update: tools/rta@${RTA_VERSION}  # updates all dependencies
	go get -u ./...
	(cd tools/optioncompare && go get -u ./...)
	go mod tidy
	go work vendor
	rm -rf tools/node_modules package-lock.json
	cd tools && ./rta npx -y npm-check-updates -u
	cd tools && ./rta npm install
	tools/rta --update
	tools/rta dprint config update
	tools/rta dprint config update --config dprint-changelog.json

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

deadcode: tools/rta@${RTA_VERSION}
	@tput bold || true
	@tput setaf 1 || true
	@tools/rta deadcode github.com/git-town/git-town/tools/format_self &
	@tools/rta deadcode github.com/git-town/git-town/tools/format_unittests &
	@tools/rta deadcode github.com/git-town/git-town/tools/stats_release &
	@tools/rta deadcode github.com/git-town/git-town/tools/structs_sorted &
	@tools/rta deadcode github.com/git-town/git-town/tools/lint_steps &
	@tools/rta deadcode -test github.com/git-town/git-town/v21 \
		| grep -v BranchExists \
		| grep -v 'Create$$' \
		| grep -v CreateFile \
		| grep -v CreateGitTown \
		| grep -v EditDefaultMessage \
		| grep -v EmptyConfigSnapshot \
		| grep -v FileExists \
		| grep -v FileHasContent \
		| grep -v IsGitRepo \
		| grep -v Memoized.AsFixture \
		| grep -v NewCommitMessages \
		| grep -v NewLineageWith \
		| grep -v NewSHAs \
		| grep -v pkg/prelude/ptr.go \
		| grep -v Paniced \
		| grep -v Set.Add \
		| grep -v UseCustomMessageOr \
		| grep -v UseDefaultMessage \
		|| true
	@tput sgr0 || true

tools/rta@${RTA_VERSION}:
	@rm -f tools/rta*
	@(cd tools && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh)
	@mv tools/rta tools/rta@${RTA_VERSION}
	@ln -s rta@${RTA_VERSION} tools/rta

tools/node_modules: tools/package-lock.json tools/rta@${RTA_VERSION}
	test -f tools/node_modules/.yarn-integrity && rm -rf tools/node_modules || true  # remove node_modules if installed with Yarn (TODO: remove after 2025-01-26)
	@echo "Installing Node based tools"
	cd tools && ./rta npm ci
	@touch tools/node_modules  # update timestamp of the node_modules folder so that Make doesn't re-install it on every command
