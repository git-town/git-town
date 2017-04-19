# How to contribute

Git Town is a welcoming community, and we'd love for everyone to bring
their contributions to make it even better.
We appreciate contributions of any size.

* Found a bug or have an idea for a new feature? - [Open an issue](https://github.com/Originate/git-town/issues/new)
* Fixed a bug or created a new feature that others will enjoy? - [Create a pull request](https://help.github.com/articles/using-pull-requests/)

This guide will help you get started and outline some things you should know when developing Git Town.


## Setup

* install [Go](https://golang.org)
  * on macOS via `brew install go`
* set up the Go directory structure on your machine
  * set the environment variable `$GOPATH` to your Go workspace
    (you can point it to any folder on your hard drive, let's assume `~/go` here)
  * add `~/go/bin` to your `$PATH`
  * create the directory `~/go/src/github.com/Originate`
  * cd into that directory, and run `git clone git@github.com:Originate/git-town.git`
  * cd into `$GOPATH/src/github.com/Originate/git-town`
  * run `bin/setup`
  * now you can run `gt` on the command line
  * see https://golang.org/doc/install#testing for details on how to test
* install [Glide](https://github.com/Masterminds/glide) (package manager for Go)
  * on macOS: `brew install glide`

* install [Ruby 2.2.3](https://www.ruby-lang.org/en/documentation/installation) to run the feature tests
  * prefer install with [rbenv](https://github.com/sstephenson/rbenv)
* install [ShellCheck](https://github.com/koalaman/shellcheck) for linting the bash scripts
* run `bundle` to install ruby gems
* optionally install [Tertestrial](https://github.com/Originate/tertestrial-server)
  for auto-running tests


## Building

* run `bin/build` to compile the source code into a runnable binary in $GOPATH/bin


## Testing

* tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).
* all features need to have comprehensive test coverage
* source code and test files must pass the linters
* See [here](./documentation/development/testing.md) for how to run the tests


## Developing

* all dependencies are located in the [vendor](vendor) folder,
  which is checked into Git
* update dependencies: `glide up`
* adding a new Go library:
  * update dependencies in a separate PR
  * `glide get <package name>`
  * your pull request for the feature that requires the new library
    should contain the updated glide files and vendor folder


## Command documentation

Every Git Town command
* has a [man page](./man/man1)
* has a [Markdown page](./documentation/commands) that is identical to the man page
* is listed on the [git-town man page](./man/man1/git-town.1)
* is listed on the [README](./README.md)


## Achitecture documents

* [branch hierarchy](./documentation/development/branch_hierarchy.md) - how Git Town sees branches
* [drivers](./documentation/development/drivers.md) - third-party specific functionality
* [steps list](./documentation/development/steps_list.md) - the architecture behind most of the Git Town commands


## Website development

* See [here](./documentation/development/website.md)


## Code style

Please follow the existing code style.
For reference, please take a look at our [Bash cheatsheet](documentation/development/bash_cheatsheet.md).


## Pull Requests

Each pull request (PR) should have the same (optional) description that it will
have when committed later and include the
[issue](https://github.com/Originate/git-town/issues) it resolves.

When merging approved PRs:
* use `git town-ship`
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
  * the version in `src/git-town` and the related features
  * the version and date in `man/man1/*.1`
* Get the feature branch reviewed and merged
* Draft a new [release](https://github.com/Originate/git-town/releases/new) against `master`

#### Homebrew/homebrew
* Fork [Homebrew](https://github.com/Homebrew/homebrew)
* Update `Library/Formula/git-town.rb`
  * Get the sha256 by downloading the release (`.tar.gz`) and using `shasum -a 256 /path/to/file`
  * Ignore the `bottle` block. It is updated by the homebrew maintainers
* Create a pull request and get it merged
