RTA_VERSION = 0.30.0  # run-that-app version to use

# internal data and state
.DEFAULT_GOAL := help
RELEASE_VERSION := "22.7.1"
GO_TEST_ARGS = LANG=C GOGC=off BROWSER=

RTA          = tools/rta@$(RTA_VERSION)
ACTIONLINT   = $(RTA) actionlint
CONC         = $(RTA) conc
CONTEST      = $(RTA) contest
CUCUMBERSORT = $(RTA) cucumber-sort
DEADCODE     = $(RTA) deadcode
DEPTH        = $(RTA) depth
DPRINT       = $(RTA) dprint
EXHAUSTRUCT  = $(RTA) exhaustruct
GHERKINLINT  = $(RTA) node node_modules/.bin/gherkin-lint
GHOKIN       = $(RTA) ghokin
GOFUMPT      = $(RTA) gofumpt
GOLANGCILINT = $(RTA) golangci-lint
NPM          = $(RTA) npm
NPX          = $(RTA) npx
NODE         = $(RTA) node
SCC          = $(RTA) scc
SHELLCHECK   = $(RTA) --optional shellcheck
SHFMT        = $(RTA) shfmt
STATICCHECK  = $(RTA) --from-source staticcheck
TAPLO        = $(RTA) taplo
TEXTRUNNER   = $(NODE) node_modules/.bin/text-runner

contest: ${RTA}  # run the Contest server
	@$(CONTEST)

cuke: install  # runs all end-to-end tests with nice output
	@env $(GO_TEST_ARGS) messyoutput=0 go test -v
	@env $(GO_TEST_ARGS) messyoutput=1 go test -v

cukeall: install  # runs all end-to-end tests on CI
	@env $(GO_TEST_ARGS) go test -v

cuke-prof: install  # creates a flamegraph for the end-to-end tests
	env $(GO_TEST_ARGS) go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

cukesmoke: install  # run the smoke E2E tests
	@env $(GO_TEST_ARGS) smoke=1 go test . -v -count=1

cukesmokewin: install  # runs the smoke E2E tests on Windows
	@env smoke=1 go test . -v -count=1

cukethis: install  # runs the end-to-end tests that have a @this tag
	@env $(GO_TEST_ARGS) cukethis=1 go test . -v -count=1

cukethiswin:  # runs the end-to-end tests that have a @this tag on Windows
	go install -ldflags "-X github.com/git-town/git-town/v22/internal/cmd.version=-dev -X github.com/git-town/git-town/v22/internal/cmd.buildDate=1/2/3"
	powershell -Command '$$env:cukethis=1 ; go test . -v -count=1'

cuke-update: install  # updates the E2E tests based on the actual behavior of Git Town
	@env $(GO_TEST_ARGS) cukeupdate=1 go test . -v -count=1
	make --no-print-directory fix

cuke-update-this: install  # updates the E2E tests that have a @this tag
	@env $(GO_TEST_ARGS) cukeupdate=1 cukethis=1 go test . -v -count=1
	make --no-print-directory fix

cukeverbose: install  # run all tests in "verbose.feature" files
	@env $(GO_TEST_ARGS) verbose=1 go test . -v -count=1

cukewin: install  # runs all end-to-end tests on Windows
	go test . -v -count=1

dependencies: ${RTA}  # prints the dependencies between the internal Go packages
	@$(DEPTH) . | grep git-town

doc: install node_modules ${RTA}  # tests the documentation
	@$(TEXTRUNNER) --offline

fix: ${RTA}  # runs all linters and auto-fixes
	make --no-print-directory fix-optioncompare-in-tests
	go run tools/format_unittests/format_unittests.go
	go run tools/format_self/format_self.go
	make --no-print-directory keep-sorted
	make --no-print-directory generate-json-schema
	$(GOFUMPT) -l -w .
	$(DPRINT) fmt
	$(DPRINT) fmt --config dprint-changelog.json
	$(SHFMT) -f . | grep -v node_modules | grep -v '^vendor/' | xargs $(SHFMT) --write
	make --no-print-directory ghokin
	tools/generate_opcodes_all.sh
	$(CUCUMBERSORT) format

generate-json-schema:  # exports the JSON-Schema for the configuration file
	(cd tools/generate_json_schema && go build) && ./tools/generate_json_schema/generate_json_schema > docs/git-town.schema.json

ghokin: ${RTA}  # formats the Cucumber tests
	@$(GHOKIN) fmt replace features/

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install:  # builds for the current platform
	@go install -ldflags="-s -w"

lint: node_modules ${RTA}  # lints the main codebase concurrently
	@$(CONC) --show=names \
		"make --no-print-directory lint-smoke" \
		"make --no-print-directory alphavet" \
		"make --no-print-directory deadcode" \
		"make --no-print-directory lint-cached-connectors" \
		"make --no-print-directory lint-collector-addf" \
		"make --no-print-directory lint-iterate-map" \
		"make --no-print-directory lint-messages-sorted" \
		"make --no-print-directory lint-messy-output" \
		"make --no-print-directory lint-optioncompare" \
		"make --no-print-directory lint-print-config" \
		"make --no-print-directory lint-structs-sorted" \
		"make --no-print-directory lint-tests-sorted" \
		"make --no-print-directory lint-use-equal" \
		"git diff --check" \
		"cd tools/lint_steps && go build && ./lint_steps" \
		"$(ACTIONLINT) " \
		"$(STATICCHECK) ./..." \
		"tools/ensure_no_files_with_dashes.sh" \
		"$(SHFMT) -f . | grep -v 'node_modules' | grep -v '^vendor/' | xargs $(SHELLCHECK)" \
		"$(GOLANGCILINT) cache clean && $(GOLANGCILINT) run" \
		"$(GHERKINLINT)" \
		"$(CUCUMBERSORT) check" \
		"make --no-print-directory lint-configfile"

lint-all: lint ${RTA}  # runs all linters
	(cd website && make test)
# $(RTA) govulncheck ./...   TODO: enable when Go 1.24.11 is available widely
	@echo lint tools/format_self
	@(cd tools/format_self && make test)
	@echo lint tools/format_unittests
	@(cd tools/format_unittests && make test)
	@echo lint tools/collector_addf
	@(cd tools/collector_addf && make test)
	@echo lint tools/generate_json_schema
	@(cd tools/generate_json_schema && make test)
	@echo lint tools/lint_cached_connectors
	@(cd tools/lint_cached_connectors && make test)
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

alphavet: ${RTA}
	@$(RTA) --available alphavet && go vet "-vettool=$(shell $(RTA) --which alphavet)" $(shell go list ./... | grep -v internal/cmd)

fix-optioncompare-in-tests:
	@(cd tools/optioncompare_in_tests && go build) && ./tools/optioncompare_in_tests/optioncompare_in_tests github.com/git-town/git-town/v22/...

keep-sorted: ${RTA}
	@$(RTA) --install ripgrep
	@$(RTA) keep-sorted $(shell $(RTA) ripgrep -l 'keep-sorted end' ./ --glob '!Makefile')

lint-cached-connectors:
	@(cd tools/lint_cached_connectors && go build) && ./tools/lint_cached_connectors/lint_cached_connectors

lint-collector-addf:
	@(cd tools/collector_addf && go build) && ./tools/collector_addf/collector_addf

lint-configfile: ${RTA}
	@$(TAPLO) check

lint-iterate-map:
	@(cd tools/iterate_map && go build) && ./tools/iterate_map/iterate_map

lint-messages-sorted:
	@(cd tools/messages_sorted && go build) && ./tools/messages_sorted/messages_sorted

lint-messy-output:
	@(cd tools/messy_output && go build) && ./tools/messy_output/messy_output

lint-print-config:
	@(cd tools/print_config_exhaustive && go build) && ./tools/print_config_exhaustive/print_config_exhaustive

lint-optioncompare:
	@(cd tools/optioncompare && go build) && ./tools/optioncompare/optioncompare github.com/git-town/git-town/v22/...

lint-smoke: ${RTA}  # runs only the essential linters
	@$(EXHAUSTRUCT) -test=false "-i=github.com/git-town/git-town.*" github.com/git-town/git-town/...
# @$(RTA) ireturn --reject="github.com/git-town/git-town/v22/pkg/prelude.Option" github.com/git-town/git-town/...

lint-structs-sorted:
	@(cd tools/structs_sorted && go build) && ./tools/structs_sorted/structs_sorted

lint-tests-sorted:
	@(cd tools/tests_sorted && go build) && ./tools/tests_sorted/tests_sorted

lint-use-equal:
	@(cd tools/use_equal && go build) && ./tools/use_equal/use_equal

stats: ${RTA}  # shows code statistics
	@find . -type f \
		| grep -v './node_modules' \
		| grep -v '\./vendor/' \
		| grep -v '\./.git/' \
		| grep -v './website/book' \
		| xargs $(SCC)

stats-release:  # displays statistics about the changes since the last release
	@(cd tools/stats_release && go build && ./stats_release v${RELEASE_VERSION})

.PHONY: test
test: install node_modules ${RTA}  # runs all the tests
	@$(CONC) --show=names \
		"make --no-print-directory cuke" \
		"make --no-print-directory doc" \
		"make --no-print-directory lint" \
		"make --no-print-directory unit"

test-go: ${RTA}  # smoke tests while working on the Go code
	@$(CONC) --show=names \
		"make --no-print-directory lint" \
		"make --no-print-directory unit"

todo:  # displays all TODO items
	@git grep --color=always --line-number TODO ':!vendor' \
		| grep -v Makefile \
		| grep -v ':= context.'

UNIT_TEST_DIRS = \
	./internal/... \
	./pkg/... \
	./tools/format_self/... \
	./tools/format_unittests/... \
	./tools/lint_cached_connectors/... \
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
	make --no-print-directory unit-text-runner

unit-text-runner: ${RTA} node_modules
	@$(NODE) --test text-runner/**/*.test.ts

unit-race: install  # runs all the unit tests with race detector
	env GOGC=off go test -count=1 -timeout 60s -race $(UNIT_TEST_DIRS)
	cd website && make --no-print-directory unit

update: ${RTA}  # updates all dependencies
	go get -u ./...
	(cd tools/optioncompare && go get -u ./...)
	go mod tidy
	go work vendor
	rm -rf node_modules package-lock.json
	$(NPX) -y npm-check-updates -u
	$(NPM) install
	$(RTA) --update
	$(DPRINT) config update
	$(DPRINT) config update --config dprint-changelog.json

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

deadcode: ${RTA}
	@$(RTA) --install deadcode
	@$(CONC) --error-on-output --show=failed \
		"$(DEADCODE) github.com/git-town/git-town/tools/format_self" \
		"$(DEADCODE) github.com/git-town/git-town/tools/format_unittests" \
		"$(DEADCODE) github.com/git-town/git-town/tools/stats_release" \
		"$(DEADCODE) github.com/git-town/git-town/tools/structs_sorted" \
		"$(DEADCODE) github.com/git-town/git-town/tools/lint_steps" \
		"$(DEADCODE) -test github.com/git-town/git-town/v22 \
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
			| grep -v NoError2 \
			| grep -v pkg/prelude/ptr.go \
			| grep -v Paniced \
			| grep -v Set.Add \
			| grep -v Set.Contains \
			| grep -v UseCustomMessageOr \
			| grep -v UseDefaultMessage \
			|| true"

tools/rta@${RTA_VERSION}:
	@rm -f tools/rta*
	@(cd tools && curl https://raw.githubusercontent.com/kevgo/run-that-app/main/download.sh | sh -s -- --version ${RTA_VERSION} --name rta@${RTA_VERSION})

node_modules: package-lock.json ${RTA}
	@echo "Installing Node based tools"
	$(NPM) ci
	@touch package-lock.json  # update timestamp so that Make doesn't re-install it on every command
	@touch node_modules  # update timestamp so that Make doesn't re-install it on every command
