.DEFAULT_GOAL := spec
date := $(shell TZ=UTC date -u '+%Y-%m-%d')

# builds for the current platform
build:
	go install -ldflags "-X github.com/Originate/git-town/src/cmd.version=v0.0.0-test -X 'github.com/Originate/git-town/src/cmd.buildDate=today'"

# builds the artifacts for a new release
build-release: cross-compile
	package/debian/make_deb.sh

# builds the binary for all platforms
cross-compile:
	go get github.com/mitchellh/gox
	gox -ldflags "-X github.com/Originate/git-town/src/cmd.version=${TRAVIS_TAG} -X 'github.com/Originate/git-town/src/cmd.buildDate=${date}'" \
			-output "dist/{{.Dir}}-{{.OS}}-{{.Arch}}"

# runs the feature tests
cuke: build
	bundle exec parallel_cucumber $(DIR)
DIR = $(if $(dir),$(dir),"features")

# deploys the website
deploy:
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

# auto-fixes lint issues in all languages
fix: fix-cucumber fix-ruby fix-markdown

# auto-fixes all Cucumber lint issues
fix-cucumber:
	bundle exec cucumber_lint --fix

# auto-fixes all Markdown lint issues
fix-markdown:
	prettier --write "{,!(vendor)/**/}*.md"

# auto-fixes all Ruby lint issues
fix-ruby:
	bundle exec rubocop --auto-correct

# lints all the source code
lint: lint-cucumber lint-go lint-markdown lint-ruby

lint-cucumber:
	bundle exec cucumber_lint

lint-go:
	goimports -d src
	gometalinter.v2

lint-markdown:
	node_modules/.bin/prettier -l '{,!(vendor)/**/}*.md'

lint-ruby:
	bundle exec rubocop

# the setup steps necessary on developer machines
setup:
	go get -u github.com/Masterminds/glide \
					  gopkg.in/alecthomas/gometalinter.v2 \
					  github.com/onsi/ginkgo/ginkgo
	gometalinter.v2 --install
	bundle install
	yarn install

# runs all the tests
spec: lint tests cuke

# runs the unit tests
tests:
	ginkgo src/...

# updates all dependencies
update:
	glide up
