VERSION ?= 0.0.0
TODAY=$(shell date +'%Y/%m/%d')

.DEFAULT_GOAL := spec

build:  # builds for the current platform
	go install -ldflags "-X github.com/git-town/git-town/src/cmd.version=v${VERSION}-dev -X github.com/git-town/git-town/src/cmd.buildDate=${TODAY}"

cuke: build   # runs the new Godog-based feature tests
	@env GOGC=off go test . -v -count=1

cuke-prof: build  # creates a flamegraph
	env GOGC=off go test . -v -cpuprofile=godog.out
	@rm git-town.test
	@echo Please open https://www.speedscope.app and load the file godog.out

docs:  # tests the documentation
	${CURDIR}/text-run/node_modules/.bin/text-run --offline

fix: fix-go fix-md  # auto-fixes lint issues in all languages

fix-go:  # auto-fixes all Go lint issues
	gofmt -s -w ./src ./test

fix-md:  # auto-fixes all Markdown lint issues
	${CURDIR}/tools/prettier/node_modules/.bin/prettier --write .

help:  # prints all make targets
	@cat Makefile | grep '^[^ ]*:' | grep -v '.PHONY' | grep -v help | sed 's/:.*#/#/' | column -s "#" -t

msi:  # compiles the MSI installer for Windows
	rm -f git-town*.msi
	go build -ldflags "-X github.com/git-town/git-town/src/cmd.version=v${VERSION} -X github.com/git-town/git-town/src/cmd.buildDate=${TODAY}"
	go-msi make --msi dist/git-town_${VERSION}_windows_intel_64.msi --version ${VERSION} --src installer/templates/ --path installer/wix.json
	@rm git-town.exe

lint: lint-go lint-md  # lints all the source code

lint-go:  # lints the Go files
	golangci-lint run src/... test/...

lint-md:   # lints the Markdown files
	${CURDIR}/tools/prettier/node_modules/.bin/prettier -l .

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
		-a dist/git-town_${VERSION}_macOS_intel_64.tar.gz \
		-a dist/git-town_${VERSION}_windows_intel_64.zip \
		v${VERSION}

release-win: msi  # adds the Windows installer to the release
	hub release edit --browse --message v${VERSION} \
		-a dist/git-town_${VERSION}_windows_intel_64.msi
		v${VERSION}

setup: setup-go  # the setup steps necessary on developer machines
	cd tools/prettier && yarn install
	cd text-run && yarn install

setup-go:
	@(cd .. && GO111MODULE=on go get github.com/cucumber/godog/cmd/godog@v0.9.0)
	@(cd .. && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0)

stats:  # shows code statistics
	@find . -type f | grep -v '\./node_modules/' | grep -v '\./vendor/' | grep -v '\./.git/' | xargs scc

test: lint docs unit cuke  # runs all the tests
.PHONY: test

test-go: build u lint-go cuke  # runs all tests for Golang

test-md: lint-md   # runs all Markdown tests

u:  # runs only the unit tests for changed code
	env GOGC=off go test -timeout 30s ./src/... ./test/...

unit:  # runs all the unit tests with race detector
	env GOGC=off go test -count=1 -timeout 60s -race ./src/... ./test/...

update:  # updates all dependencies
	go get -u ./...
	go mod tidy
	go mod vendor

website-build:  # compiles the website (used during deployment)
	(cd tools/harp && yarn install)
	tools/harp/node_modules/.bin/harp compile website/ www

website-dev:  # runs a local development server of the website
	(cd website && ../tools/harp/node_modules/.bin/harp server)
