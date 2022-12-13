VERSION ?= 0.0.0
TODAY=$(shell date +'%Y/%m/%d')
.DEFAULT_GOAL := help

DEPTH_VERSION = 1.2.1
GOFUMPT_VERSION = 0.3.0
GOLANGCILINT_VERSION = 1.50.0
SCC_VERSION = 3.1.0
SHELLCHECK_VERSION = 0.8.0
SHFMT_VERSION = 3.5.1

build:  # builds for the current platform
	go install -trimpath -ldflags "-X github.com/git-town/git-town/v7/src/cmd.version=v${VERSION}-dev -X github.com/git-town/git-town/v7/src/cmd.buildDate=${TODAY}"

cuke: build   # runs all end-to-end tests
	@env LANG=C GOGC=off go test . -v -count=1

cukethis: build   # runs the end-to-end tests that have a @this tag
	@env LANG=C GOGC=off go test . -v -count=1 -this

cuke-prof: build  # creates a flamegraph
	env LANG=C GOGC=off go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

dependencies: tools/depth-${DEPTH_VERSION}  # prints the dependencies between packages as a tree
	@tools/depth-${DEPTH_VERSION} . | grep git-town

docs: build tools/node_modules  # tests the documentation
	${CURDIR}/tools/node_modules/.bin/text-run --offline

fix: tools/golangci-lint-${GOLANGCILINT_VERSION} tools/gofumpt-${GOFUMPT_VERSION} tools/node_modules tools/shellcheck-${SHELLCHECK_VERSION} tools/shfmt-${SHFMT_VERSION}  # auto-fixes lint issues in all languages
	git diff --check
	tools/gofumpt-${GOFUMPT_VERSION} -l -w .
	${CURDIR}/tools/node_modules/.bin/dprint fmt
	${CURDIR}/tools/node_modules/.bin/prettier --write '**/*.yml'
	tools/shfmt-${SHFMT_VERSION} -f . | grep -v tools/node_modules | grep -v '^vendor\/' | xargs tools/shfmt-${SHFMT_VERSION} --write
	tools/shfmt-${SHFMT_VERSION} -f . | grep -v tools/node_modules | grep -v '^vendor\/' | xargs tools/shellcheck-${SHELLCHECK_VERSION}
	${CURDIR}/tools/node_modules/.bin/gherkin-lint
	tools/golangci-lint-${GOLANGCILINT_VERSION} run

help:  # prints all available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v '.PHONY' | grep -v help | grep -v "^tools\/" | sed 's/:.*#/#/' | column -s "#" -t

msi:  # compiles the MSI installer for Windows
	rm -f git-town*.msi
	go build -trimpath -ldflags "-X github.com/git-town/git-town/src/cmd.version=v${VERSION} -X github.com/git-town/git-town/src/cmd.buildDate=${TODAY}"
	go-msi make --msi dist/git-town_${VERSION}_windows_intel_64.msi --version ${VERSION} --src installer/templates/ --path installer/wix.json
	@rm git-town.exe

release-linux:   # creates a new release
	# cross-compile the binaries
	goreleaser --rm-dist

	# create GitHub release with files in alphabetical order
	hub release create --draft --browse --message v${VERSION} \
		-a dist/git-town_${VERSION}_linux_intel_64.deb \
		-a dist/git-town_${VERSION}_linux_intel_64.rpm \
		-a dist/git-town_${VERSION}_linux_intel_64.tar.gz \
		-a dist/git-town_${VERSION}_linux_arm_64.deb \
		-a dist/git-town_${VERSION}_linux_arm_64.rpm \
		-a dist/git-town_${VERSION}_linux_arm_64.tar.gz \
		-a dist/git-town_${VERSION}_macos_intel_64.tar.gz \
		-a dist/git-town_${VERSION}_macos_arm_64.tar.gz \
		-a dist/git-town_${VERSION}_windows_intel_64.zip \
		v${VERSION}

release-win: msi  # adds the Windows installer to the release
	hub release edit --browse --message v${VERSION} \
		-a dist/git-town_${VERSION}_windows_intel_64.msi
		v${VERSION}

stats: tools/scc-${SCC_VERSION}  # shows code statistics
	@find . -type f | grep -v './tools/node_modules' | grep -v '\./vendor/' | grep -v '\./.git/' | grep -v './website/book' | xargs tools/scc-${SCC_VERSION}

test: fix docs unit cuke  # runs all the tests
.PHONY: test

unit:  # runs only the unit tests for changed code
	env GOGC=off go test -timeout 30s ./src/... ./test/...

unit-all:  # runs all the unit tests with race detector
	env GOGC=off go test -count=1 -timeout 60s -race ./src/... ./test/...

update:  # updates all dependencies
	go get -u ./...
	go mod tidy
	go mod vendor
	(cd tools && yarn upgrade --latest)
	echo
	echo Please update the third-party tooling in the Makefile manually.

# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

tools/depth-${DEPTH_VERSION}:
	@echo "Installing depth ${DEPTH_VERSION} ..."
	@env GOBIN="$(CURDIR)/tools" go install github.com/KyleBanks/depth/cmd/depth@v${DEPTH_VERSION}
	@mv tools/depth tools/depth-${DEPTH_VERSION}

tools/gofumpt-${GOFUMPT_VERSION}:
	@echo "Installing gofumpt ${GOFUMPT_VERSION} ..."
	@env GOBIN="$(CURDIR)/tools" go install mvdan.cc/gofumpt@v${GOFUMPT_VERSION}
	@mv tools/gofumpt tools/gofumpt-${GOFUMPT_VERSION}

tools/golangci-lint-${GOLANGCILINT_VERSION}:
	@echo "Installing golangci-lint ${GOLANGCILINT_VERSION} ..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b tools v${GOLANGCILINT_VERSION}
	@mv tools/golangci-lint tools/golangci-lint-${GOLANGCILINT_VERSION}

tools/node_modules: tools/yarn.lock
	@echo "Installing Node based tools"
	@cd tools && yarn install
	@touch tools/node_modules  # update timestamp of the node_modules folder so that Make doesn't re-install it on every command

tools/scc-${SCC_VERSION}:
	@echo "Installing scc ${SCC_VERSION} ..."
	@env GOBIN=$(CURDIR)/tools go install github.com/boyter/scc/v3@v3.1.0
	@mv tools/scc tools/scc-${SCC_VERSION}

tools/shellcheck-${SHELLCHECK_VERSION}:
	@echo installing Shellcheck ${SHELLCHECK_VERSION} ...
	@curl -sSL https://github.com/koalaman/shellcheck/releases/download/v${SHELLCHECK_VERSION}/shellcheck-v${SHELLCHECK_VERSION}.linux.x86_64.tar.xz | tar xJ
	@mv shellcheck-v${SHELLCHECK_VERSION}/shellcheck tools/shellcheck-${SHELLCHECK_VERSION}
	@rm -rf shellcheck-v${SHELLCHECK_VERSION}

tools/shfmt-${SHFMT_VERSION}:
	@echo installing Shellfmt ${SHFMT_VERSION} ...
	@env GOBIN="$(CURDIR)/tools" go install mvdan.cc/sh/v3/cmd/shfmt@v${SHFMT_VERSION}
	@mv tools/shfmt tools/shfmt-${SHFMT_VERSION}
