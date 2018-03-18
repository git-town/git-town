.DEFAULT_GOAL := spec


# Builds for the current platform
build:
	go install

# Makes a new binary release
build-release: cross-compile
	package/debian/make_deb.sh

# Builds the binary for all platforms
cross-compile:
	go get github.com/mitchellh/gox
	timestamp=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')
	sha=$(git rev-parse HEAD)
	gox -ldflags "-X github.com/Originate/git-town/cmd.Version=$TRAVIS_TAG -X github.com/Originate/git-town/cmd.BuildTime=$timestamp) -X github.com/Originate/git-town/cmd.GitHash=$sha" \
			-output "dist/{{.Dir}}-{{.OS}}-{{.Arch}}"

# Runs the feature tests
cuke: build
	bundle exec parallel_cucumber $(DIR)
DIR = $(if $(dir),$(dir),"features")

# Deploys the website
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

# Fixes all issues in all languages
fix: fix-cucumber fix-ruby fix-markdown

# Fixes all Cucumber issues
fix-cucumber:
	bundle exec cucumber_lint --fix

# Fixes all Markdown issues
fix-markdown:
	prettier --write "{,!(vendor)/**/}*.md"

# Fixes all Ruby issues
fix-ruby:
	bundle exec rubocop --auto-correct

# Lints all the source code
lint: lint-go lint-markdown lint-ruby

lint-go:
	goimports -d src
	gometalinter.v2

lint-markdown:
	node_modules/.bin/prettier -l '{,!(vendor)/**/}*.md'

lint-ruby:
	bundle exec rubocop

# The setup steps necessary on developer machines
setup:
	go get github.com/Masterminds/glide \
				 github.com/onsi/ginkgo/ginkgo
	go get -u gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install
	bundle install
	yarn install

# Runs all the tests
spec: lint tests cuke

# Runs the unit tests
tests:
	ginkgo src/...

# Updates all dependencies
update:
	glide up
