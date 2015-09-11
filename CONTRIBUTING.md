# How to contribute

Git Town is a welcoming community, and we'd love for everyone to bring
their contributions to make it even better.
We appreciate contributions of any size.

* Found a bug or have an idea for a new feature? - [Open an issue](https://github.com/Originate/git-town/issues/new)
* Fixed a bug or created a new feature that others will enjoy? - [Create a pull request](https://help.github.com/articles/using-pull-requests/)

This guide will help you get started and outline some things you should know when developing Git Town.


## Setup

* fork and clone the repository to your machine
* install [Ruby 2.2.3](https://www.ruby-lang.org/en/documentation/installation) to run the feature tests
  * prefer install with [rbenv](https://github.com/sstephenson/rbenv)
* install [ShellCheck](https://github.com/koalaman/shellcheck) for linting the bash scripts
* run `bundle` to install ruby gems


## Testing

* tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).
* all features need to have comprehensive test coverage
* source code and test files must pass the linters
* See [here](./documentation/development/testing.md) for how to run the tests


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


## Pull Requests

Each pull request (PR) should have the same (optional) description that it will
have when committed later and include the
[issue](https://github.com/Originate/git-town/issues) it resolves.

When merging approved PRs:
* use `git ship`
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
