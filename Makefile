VERSION ?= 0.0.0
TODAY=$(shell date +'%Y/%m/%d')
.DEFAULT_GOAL := help

build:  # builds for the current platform
	go install -ldflags "-X github.com/git-town/git-town/v7/src/cmd.version=v${VERSION}-dev -X github.com/git-town/git-town/v7/src/cmd.buildDate=${TODAY}"

cuke: build   # runs all end-to-end tests
	@env LANG=C GOGC=off go test . -v -count=1

cukethis: build   # runs the end-to-end tests that have a @this tag
	@env LANG=C GOGC=off go test . -v -count=1 -this

cuke-prof: build  # creates a flamegraph
	env LANG=C GOGC=off go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

dependencies: tools/depth  # prints the dependencies between packages as a tree
	@tools/depth . | grep git-town

docs: build tools/node_modules  # tests the documentation
	${CURDIR}/tools/node_modules/.bin/text-run --offline

fix: tools/gofumpt  # auto-fixes lint issues in all languages
	tools/gofumpt -l -w .
	dprint fmt
	${CURDIR}/tools/node_modules/.bin/prettier --write '**/*.yml'

help:  # prints all available targets
	@cat Makefile | grep '^[^ ]*:' | grep -v '.PHONY' | grep -v help | grep -v "^tools\/" | sed 's/:.*#/#/' | column -s "#" -t

lint: tools/golangci-lint tools/node_modules  # lints all the source code
	git diff --check
	tools/golangci-lint run
	${CURDIR}/tools/node_modules/.bin/dprint check
	${CURDIR}/tools/node_modules/.bin/prettier --check '**/*.yml'

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

stats: tools/scc  # shows code statistics
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
	(cd tools && yarn upgrade --latest)
	echo
	echo Please update the third-party tooling in the Makefile manually.


# --- HELPER TARGETS --------------------------------------------------------------------------------------------------------------------------------

tools/depth: Makefile
	env GOBIN="$(CURDIR)/tools" go install github.com/KyleBanks/depth/cmd/depth@latest

tools/gofumpt: Makefile
	env GOBIN="$(CURDIR)/tools" go install mvdan.cc/gofumpt@v0.3.0

tools/golangci-lint: Makefile
	@echo "Installing golangci-lint ..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b tools v1.50.0

tools/node_modules: tools/yarn.lock
	@cd tools && yarn install
	@touch tools/node_modules  # update timestamp of the node_modules folder so that Make doesn't re-install it on every command

tools/scc: Makefile
	env GOBIN="$(CURDIR)/tools" go install github.com/boyter/scc@latest
