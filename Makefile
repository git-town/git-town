.DEFAULT_GOAL := spec


build:  # builds for the current platform
	go install

build-release: cross-compile  # builds the artifacts for a new release
	package/debian/make_deb.sh

cross-compile:  # builds the binary for all platforms
	go get github.com/mitchellh/gox
	timestamp=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')
	sha=$(git rev-parse HEAD)
	gox -ldflags "-X github.com/Originate/git-town/cmd.Version=$TRAVIS_TAG -X github.com/Originate/git-town/cmd.BuildTime=$timestamp) -X github.com/Originate/git-town/cmd.GitHash=$sha" \
			-output "dist/{{.Dir}}-{{.OS}}-{{.Arch}}"

cuke: build  # runs the feature tests
	bundle exec parallel_cucumber $(DIR)
DIR = $(if $(dir),$(dir),"features")

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

fix: fix-cucumber fix-ruby fix-markdown  # auto-fixes lint issues in all languages

fix-cucumber:  # auto-fixes all Cucumber lint issues
	bundle exec cucumber_lint --fix

fix-markdown:  # auto-fixes all Markdown lint issues
	prettier --write "{,!(vendor)/**/}*.md"

fix-ruby:  # auto-fixes all Ruby lint issues
	bundle exec rubocop --auto-correct

help:  # prints all make targets
	@cat Makefile | grep '^[^ ]*:' | grep -v '.PHONY' | grep -v help | sed 's/:.*#/#/' | column -s "#" -t

lint: lint-cucumber lint-go lint-markdown lint-ruby  # lints all the source code

lint-cucumber:  # lints the Cucumber files
	bundle exec cucumber_lint

lint-go:  # lints the Go files
	goimports -d src
	gometalinter.v2

lint-markdown:  # lints the Markdown files
	node_modules/.bin/prettier -l '{,!(vendor)/**/}*.md'
	node_modules/.bin/text-run --offline

lint-ruby:  # lints the Ruby files
	bundle exec rubocop

setup:  # the setup steps necessary on developer machines
	go get -u github.com/Masterminds/glide \
					  gopkg.in/alecthomas/gometalinter.v2 \
					  github.com/onsi/ginkgo/ginkgo
	gometalinter.v2 --install
	bundle install
	yarn install

spec: lint tests cuke  # runs all the tests

tests:  # runs the unit tests
	ginkgo src/...

update:  # updates all dependencies
	glide up
