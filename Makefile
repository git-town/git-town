# dev tooling and versions
DEPTH_VERSION = 1.2.1
GOFUMPT_VERSION = 0.4.0
GOLANGCILINT_VERSION = 1.50.1
SCC_VERSION = 3.1.0
SHELLCHECK_VERSION = 0.9.0
SHFMT_VERSION = 3.6.0

# internal data and state
.DEFAULT_GOAL := help
TODAY = $(shell date +'%Y-%m-%d')
DEV_VERSION := $(shell git describe --tags 2> /dev/null || git rev-parse --short HEAD)
RELEASE_VERSION := "9.0.0"
GO_BUILD_ARGS = LANG=C GOGC=off

build:  # builds for the current platform
	go install -ldflags "-X github.com/git-town/git-town/v9/src/cmd.version=${DEV_VERSION}-dev -X github.com/git-town/git-town/v9/src/cmd.buildDate=${TODAY}"

cuke: build   # runs all end-to-end tests
	@env $(GO_BUILD_ARGS) go test . -v -count=1

cukethis: build   # runs the end-to-end tests that have a @this tag
	@env $(GO_BUILD_ARGS) cukethis=1 go test . -v -count=1

cuke-prof: build  # creates a flamegraph for the end-to-end tests
	env $(GO_BUILD_ARGS) go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

dependencies: tools/depth_${DEPTH_VERSION}  # prints the dependencies between the internal Go packages as a tree
	@tools/depth_${DEPTH_VERSION} . | grep git-town

docs: build tools/node_modules  # tests the documentation
	${CURDIR}/tools/node_modules/.bin/text-run --offline

fix: tools/golangci_lint_${GOLANGCILINT_VERSION} tools/gofumpt_${GOFUMPT_VERSION} tools/node_modules tools/shellcheck_${SHELLCHECK_VERSION} tools/shfmt_${SHFMT_VERSION}  # auto-fixes lint issues in all languages
	git diff --check
	tools/gofumpt_${GOFUMPT_VERSION} -l -w .
	${CURDIR}/tools/node_modules/.bin/dprint fmt
	${CURDIR}/tools/node_modules/.bin/prettier --write '**/*.yml'
	tools/shfmt_${SHFMT_VERSION} -f . | grep -v tools/node_modules | grep -v '^vendor/' | xargs tools/shfmt_${SHFMT_VERSION} --write
	tools/shfmt_${SHFMT_VERSION} -f . | grep -v tools/node_modules | grep -v '^vendor/' | xargs tools/shellcheck_${SHELLCHECK_VERSION}
	${CURDIR}/tools/node_modules/.bin/gherkin-lint
	tools/golangci_lint_${GOLANGCILINT_VERSION} run
	tools/ensure_no_files_with_dashes.sh

help:  # prints all available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

msi: version_tag_is_up_to_date  # compiles the MSI installer for Windows
	rm -f git-town*.msi
	go build -trimpath -ldflags "-X github.com/git-town/git-town/src/cmd.version=${RELEASE_VERSION} -X github.com/git-town/git-town/src/cmd.buildDate=${TODAY}"
	go-msi make --msi dist/git-town_${RELEASE_VERSION}_windows_intel_64.msi --version ${RELEASE_VERSION} --src installer/templates/ --path installer/wix.json
	@rm git-town.exe

release-linux: version_tag_is_up_to_date   # creates a new release
	# cross-compile the binaries
	goreleaser --rm-dist
	# create GitHub release with files in alphabetical order
	hub release create --draft --browse --message "v${RELEASE_VERSION}" \
		-a dist/git-town_${RELEASE_VERSION}_linux_intel_64.deb \
		-a dist/git-town_${RELEASE_VERSION}_linux_intel_64.rpm \
		-a dist/git-town_${RELEASE_VERSION}_linux_intel_64.tar.gz \
		-a dist/git-town_${RELEASE_VERSION}_linux_arm_64.deb \
		-a dist/git-town_${RELEASE_VERSION}_linux_arm_64.rpm \
		-a dist/git-town_${RELEASE_VERSION}_linux_arm_64.tar.gz \
		-a dist/git-town_${RELEASE_VERSION}_macos_intel_64.tar.gz \
		-a dist/git-town_${RELEASE_VERSION}_macos_arm_64.tar.gz \
		-a dist/git-town_${RELEASE_VERSION}_windows_intel_64.zip \
		"v${RELEASE_VERSION}"

release-win: msi version_tag_is_up_to_date  # adds the Windows installer to the release
	hub release edit \
		-a dist/git-town_${RELEASE_VERSION}_windows_intel_64.msi \
		v${RELEASE_VERSION}

stats: tools/scc_${SCC_VERSION}  # shows code statistics
	@find . -type f | grep -v './tools/node_modules' | grep -v '\./vendor/' | grep -v '\./.git/' | grep -v './website/book' | xargs tools/scc_${SCC_VERSION}

test: fix docs unit cuke  # runs all the tests
.PHONY: test

test-go: tools/gofumpt_${GOFUMPT_VERSION} tools/golangci_lint_${GOLANGCILINT_VERSION}  # smoke tests for Go refactorings
	tools/gofumpt_${GOFUMPT_VERSION} -l -w . &
	make --no-print-directory unit &
	make --no-print-directory build &
	tools/golangci_lint_${GOLANGCILINT_VERSION} run

todo:  # displays all TODO items
	git grep --line-number -C1 TODO ':!vendor'

unit:  # runs only the unit tests for changed code
	env GOGC=off go test -timeout 30s ./src/... ./test/...

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

tools/depth_${DEPTH_VERSION}:
	@echo "Installing depth ${DEPTH_VERSION} ..."
	@env GOBIN="$(CURDIR)/tools" go install github.com/KyleBanks/depth/cmd/depth@v${DEPTH_VERSION}
	@mv tools/depth tools/depth_${DEPTH_VERSION}

tools/gofumpt_${GOFUMPT_VERSION}:
	@echo "Installing gofumpt ${GOFUMPT_VERSION} ..."
	@env GOBIN="$(CURDIR)/tools" go install mvdan.cc/gofumpt@v${GOFUMPT_VERSION}
	@mv tools/gofumpt tools/gofumpt_${GOFUMPT_VERSION}

tools/golangci_lint_${GOLANGCILINT_VERSION}:
	@echo "Installing golangci-lint ${GOLANGCILINT_VERSION} ..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b tools v${GOLANGCILINT_VERSION}
	@mv tools/golangci-lint tools/golangci_lint_${GOLANGCILINT_VERSION}

tools/node_modules: tools/yarn.lock
	@echo "Installing Node based tools"
	@cd tools && yarn install
	@touch tools/node_modules  # update timestamp of the node_modules folder so that Make doesn't re-install it on every command

tools/scc_${SCC_VERSION}:
	@echo "Installing scc ${SCC_VERSION} ..."
	@env GOBIN=$(CURDIR)/tools go install github.com/boyter/scc/v3@v3.1.0
	@mv tools/scc tools/scc_${SCC_VERSION}

tools/shellcheck_${SHELLCHECK_VERSION}:
	@echo installing Shellcheck ${SHELLCHECK_VERSION} ...
	@curl -sSL https://github.com/koalaman/shellcheck/releases/download/v${SHELLCHECK_VERSION}/shellcheck-v${SHELLCHECK_VERSION}.$(shell go env GOOS).x86_64.tar.xz | tar xJ
	@mv shellcheck-v${SHELLCHECK_VERSION}/shellcheck tools/shellcheck_${SHELLCHECK_VERSION}
	@rm -rf shellcheck-v${SHELLCHECK_VERSION}

tools/shfmt_${SHFMT_VERSION}:
	@echo installing Shellfmt ${SHFMT_VERSION} ...
	@env GOBIN="$(CURDIR)/tools" go install mvdan.cc/sh/v3/cmd/shfmt@v${SHFMT_VERSION}
	@mv tools/shfmt tools/shfmt_${SHFMT_VERSION}

# verifies that the latest commit in the repo has a Git tag
version_tag_is_up_to_date:
	@[ ! -z "$(RELEASE_VERSION)" ] || (echo "Please add an up-to-date Git tag for the release"; exit 5)
