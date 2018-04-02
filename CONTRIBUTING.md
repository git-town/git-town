# How to contribute

Git Town is a welcoming community,
and we'd love for everyone to bring
their contributions of any size to make it even better.

* Found a bug or have an idea for a new feature? - [Open an issue](https://github.com/Originate/git-town/issues/new)
* Fixed a bug or created a new feature that others will enjoy? - [Create a pull request](https://help.github.com/articles/using-pull-requests/)

This guide will help you get started and outline some things you should know when developing Git Town.

## Setup

* install [Go](https://golang.org) version 1.9 or higher
  * on macOS via `brew install go`
  * on Windows via the [official Go installer](https://golang.org/dl)
* install [Ruby 2.4.1](https://www.ruby-lang.org/en/documentation/installation) to run the feature tests
  * prefer install with [rbenv](https://github.com/sstephenson/rbenv)
  * run `gem install bundler`
* set up the Go directory structure on your machine
  * set the environment variable `$GOPATH` to your Go workspace
    (you can point it to any folder on your hard drive, let's assume `~/go` here)
  * add `~/go/bin` to your `$PATH`
  * create the directory `~/go/src/github.com/Originate`
  * cd into that directory, and run `git clone git@github.com:Originate/git-town.git`
  * cd into `$GOPATH/src/github.com/Originate/git-town`
* make sure you have `make` - Mac and Linux users should be okay,
  Windows users should install
  [Make for Windows](http://gnuwin32.sourceforge.net/packages/make.htm)
* run `make setup` and then `make build`
* now you can run `git-town` on the command line
* see https://golang.org/doc/install#testing for details on how to test
* optionally install [Tertestrial](https://github.com/Originate/tertestrial-server)
  for auto-running tests

## Building

* run `make build` to compile the source code into a runnable binary in $GOPATH/bin

## Testing

* tests are written in [Cucumber](http://cukes.info) and [RSpec](http://rspec.info).
* all features need to have comprehensive test coverage
* source code and test files must pass the linters
* See [here](./documentation/development/testing.md) for how to run the tests

## Developing

* all dependencies are located in the [vendor](vendor) folder,
  which is checked into Git
* update dependencies: `make update`
* adding a new Go library:
  * `glide get <package name>`
  * your pull request for the feature that requires the new library
    should contain the updated glide files and vendor folder

## Command documentation

Every Git Town command

* has a [Markdown page](./documentation/commands) that is identical to the man page
* is listed on the [README](./README.md)

## Achitecture documents

* [branch hierarchy](./documentation/development/branch_hierarchy.md) - how Git Town sees branches
* [drivers](./documentation/development/drivers.md) - third-party specific functionality
* [steps list](./documentation/development/steps_list.md) - the architecture behind most of the Git Town commands

## Website development

* See [here](./documentation/development/website.md)

## Code style

Please follow the existing code style.

## Pull Requests

Each pull request (PR) should have the same (optional) description that it will
have when committed later and include the
[issue](https://github.com/Originate/git-town/issues) it resolves.

When merging approved PRs:

* use `git town ship`
* the message for the squashed commit should follow the
  [formatting guidelines for commit messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
* [mention the closed issue](https://help.github.com/articles/closing-issues-via-commit-messages)
  in the commit body, so that the respective issue is automatically closed.

Example of a commit message:

```
Automatically foo commits in "git bar"

Fooing changes before they get barred into a separate branch
keeps the final foo simpler when shipping that branch later.

Implements #123
```

## Release Process

#### Originate/git-town

* Create a feature branch which updates
  * `RELEASE_NOTES.md`
  * the version in `src/cmd/version.go` and the related features
* Get the feature branch reviewed and merged
* Create and push a new Git Tag for the release
  * `git tag -m release -a v4.0`
  * `git push --tags`
* Travis-CI creates a new release on Github and attaches the GT binaries to it

#### Homebrew/homebrew

* Fork [Homebrew](https://github.com/Homebrew/homebrew)
* Update `Library/Formula/git-town.rb`
  * Get the sha256 by downloading the release (`.tar.gz`) and using `shasum -a 256 /path/to/file`
  * Ignore the `bottle` block. It is updated by the homebrew maintainers
* Create a pull request and get it merged
