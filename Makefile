.DEFAULT_GOAL := spec


# builds for the current platform
build:
	go install

# makes a new binary release
build-release: cross-compile
	package/debian/make_deb.sh

# builds the binary for all platforms
cross-compile:
	go get github.com/mitchellh/gox
	timestamp=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')
	sha=$(git rev-parse HEAD)
	gox -ldflags "-X github.com/Originate/git-town/cmd.Version=$TRAVIS_TAG -X github.com/Originate/git-town/cmd.BuildTime=$timestamp) -X github.com/Originate/git-town/cmd.GitHash=$sha" \
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

# fixes all issues in all languages
fix: fix-cucumber fix-ruby fix-markdown

# fixes all Cucumber issues
fix-cucumber:
	bundle exec cucumber_lint --fix

# fixes all Markdown issues
fix-markdown:
	prettier --write "{,!(vendor)/**/}*.md"

# fixes all Ruby issues
fix-ruby:
	bundle exec rubocop --auto-correct

# lints all the source code
lint: lint-go lint-markdown lint-ruby

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
