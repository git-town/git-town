.DEFAULT_GOAL := spec
date := $(shell TZ=UTC date -u '+%Y-%m-%d')

build:  # builds for the current platform
	go install -ldflags "-X github.com/git-town/git-town/src/cmd.version=v0.0.0-test -X github.com/git-town/git-town/src/cmd.buildDate=today"

build-release: cross-compile  # builds the artifacts for a new release
	package/debian/make_deb.sh

cross-compile:  # builds the binary for all platforms
	go get github.com/mitchellh/gox
	gox -ldflags "-X github.com/git-town/git-town/src/cmd.version=${TRAVIS_TAG} -X github.com/git-town/git-town/src/cmd.buildDate=${date}" \
			-output "dist/{{.Dir}}-${TRAVIS_TAG}-{{.OS}}-{{.Arch}}"

cuke: cuke-go cuke-rb  # runs the feature tests

cuke-go: build   # runs the new Godog-based feature tests
	godog --concurrency=$(shell nproc --all) --format=progress features/git-town features/git-town-alias features/git-town-append features/git-town-config features/git-town-hack/on_feature_branch/in_committed_subfolder.feature features/git-town-hack/on_feature_branch/with_remote_origin.feature

cuke-rb: build   # runs the old Ruby-based feature tests
	bundle exec parallel_cucumber features/git-town-hack/branch_exists_locally.feature features/git-town-hack/branch_exists_remotely.feature features/git-town-hack/branch_name_with_slash.feature features/git-town-hack/hack_push_flag.feature features/git-town-hack/main_branch_rebase_tracking_branch_conflict.feature features/git-town-hack/offline.feature features/git-town-hack/on_feature_branch/in_uncommitted_subfolder.feature features/git-town-hack/on_feature_branch/without_remote_origin.feature features/git-town-hack/on_main_branch/with_upstream.feature features/git-town-hack/on_main_branch/without_remote_origin.feature features/git-town-hack/on_main_branch/in_subfolder.feature features/git-town-hack/prompt_for_parent.feature features/git-town-install-fish-autocompletion features/git-town-kill features/git-town-main_branch features/git-town-new-branch-push-flag features/git-town-new-pull-request features/git-town-offline-mode features/git-town-perennial_branches features/git-town-prepend features/git-town-prune-branches features/git-town-pull_branch_strategy features/git-town-rename-branch features/git-town-repo features/git-town-set-parent-branch features/git-town-ship features/git-town-sync features/git-town-version

deploy:  # deploys the website
	git checkout gh-pages
	git pull
	git checkout master
	git pull --rebase
	harp compile website/ _www
	git checkout gh-pages
	cp -r _www/* .
	rm -rf _www
	git add -A
	git commit
	git push
	git checkout master

fix: fix-cucumber fix-go fix-rb fix-md  # auto-fixes lint issues in all languages

fix-cucumber:  # auto-fixes all Cucumber lint issues
	bundle exec cucumber_lint --fix

fix-go:  # auto-fixes all Go lint issues
	gofmt -s -w ./src ./test

fix-md:  # auto-fixes all Markdown lint issues
	@find . -type f \( \
		-path '**/*.md' -o \
		-path '**/*.yml' -o \
		-path '**/*.json' -o \
		-path '**/*.js' \) | \
		grep -v node_modules | \
		grep -v vendor | \
		xargs node_modules/.bin/prettier --write

fix-rb:  # auto-fixes all Ruby lint issues
	bundle exec rubocop --auto-correct

help:  # prints all make targets
	@cat Makefile | grep '^[^ ]*:' | grep -v '.PHONY' | grep -v help | sed 's/:.*#/#/' | column -s "#" -t

lint: lint-cucumber lint-go lint-md lint-rb  # lints all the source code

lint-cucumber:  # lints the Cucumber files
	bundle exec cucumber_lint

lint-go:  # lints the Go files
	golangci-lint run --enable-all -D dupl -D lll -D gochecknoglobals -D gochecknoinits -D goconst -D wsl -D gomnd src/... test/...

lint-md: build  # lints the Markdown files
	@find . -type f \( \
		-path '**/*.md' -o \
		-path '**/*.yml' -o \
		-path '**/*.json' -o \
		-path '**/*.js' \) | \
		grep -v node_modules | \
		grep -v vendor | \
		xargs node_modules/.bin/prettier -l
	node_modules/.bin/text-run --offline

lint-rb:  # lints the Ruby files
	bundle exec rubocop

setup:  # the setup steps necessary on developer machines
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	GO111MODULE=on go get github.com/cucumber/godog/cmd/godog@v0.8.1
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(shell go env GOPATH)/bin v1.23.8
	bundle install
	yarn install

stats:  # shows code statistics
	@find . -type f | grep -v '\./node_modules/' | grep -v '\./vendor/' | grep -v '\./.git/' | xargs scc

test: lint unit cuke  # runs all the tests
.PHONY: test

test-go: build unit cuke-go lint-go  # runs all tests for Golang

u:  # runs only the unit tests for changed code
	go test -timeout 3s ./src/... ./test/...

unit:  # runs all the unit tests with race detector
	go test -count=1 -timeout 3s -race ./src/... ./test/...

update:  # updates all dependencies
	dep ensure -update
