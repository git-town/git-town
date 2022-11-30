VERSION ?= 0.0.0
TODAY=$(shell date +'%Y/%m/%d')

build:  # builds for the current platform
	go install -ldflags "-X github.com/git-town/git-town/v7/src/cmd.version=v${VERSION}-dev -X github.com/git-town/git-town/v7/src/cmd.buildDate=${TODAY}"

cuke: build   # runs all end-to-end tests
	@env LANG=C GOGC=off go test . -v -count=1

cuke-open:  # runs only the currently uncommitted end-to-end tests
	@git status --porcelain | grep -v '^\s*D ' | sed 's/^\s*\w\s*//' | grep '\.feature' | xargs godog
#                           remove deleted     remove indicator

cukethis: build   # runs the end-to-end tests that have a @this tag
	@env LANG=C GOGC=off go test . -v -count=1 -this

cuke-prof: build  # creates a flamegraph
	env LANG=C GOGC=off go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

dependencies:  # prints the dependencies between packages as a tree
	@depth . | grep git-town

docs: build  # tests the documentation
	${CURDIR}/tools/node_modules/.bin/text-run --offline

fix: fix-go fix-md  # auto-fixes lint issues in all languages

fix-go: tools/gofumpt  # auto-fixes all Go lint issues
	tools/gofumpt -l -w .

fix-md:  # auto-fixes all Markdown lint issues
	dprint fmt

help:  # prints all available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v '.PHONY' | grep -v help | sed 's/:.*#/#/' | column -s "#" -t

lint: lint-go lint-md  # lints all the source code
	git diff --check

lint-go: tools/golangci-lint  # lints the Go files
	tools/golangci-lint run

lint-md:   # lints the Markdown files
	dprint check

msi:  # compiles the MSI installer for Windows
	rm -f git-town*.msi
	go build -ldflags "-X github.com/git-town/git-town/src/cmd.version=v${VERSION} -X github.com/git-town/git-town/src/cmd.buildDate=${TODAY}"
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

setup: setup-go setup-tools  # the setup steps necessary on developer machines

setup-tools:  # the setup steps necessary for document tests
	cd tools && yarn install

setup-go: setup-godog
	go install github.com/KyleBanks/depth/cmd/depth@latest
	go install github.com/boyter/scc@latest

setup-godog:  # install the godog binary
	go install github.com/cucumber/godog/cmd/godog@v0.9.0

stats:  # shows code statistics
	@find . -type f | grep -v './tools/node_modules' | grep -v '\./vendor/' | grep -v '\./.git/' | grep -v './website/book' | xargs scc

test: lint docs u cuke  # runs all the tests
.PHONY: test

test-go: build u lint-go cuke  # runs all tests for Golang

test-md: lint-md   # runs all Markdown tests

unit:  # runs only the unit tests for changed code
	env GOGC=off go test -timeout 30s ./src/... ./test/...

unit-all:  # runs all the unit tests with race detector
	env GOGC=off go test -count=1 -timeout 60s -race ./src/... ./test/...

update:  # updates all dependencies
	go get -u ./...
	go mod tidy
	go mod vendor
	echo
	echo Please update the tools that "make setup" installs manually.


# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

tools/gofumpt: Makefile
	env GOBIN="$(CURDIR)/tools" go install mvdan.cc/gofumpt@v0.3.0

tools/golangci-lint: Makefile
	@echo "Installing golangci-lint ..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b tools v1.50.0



.DEFAULT_GOAL := help
