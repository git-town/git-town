# How to contribute

Git Town is a welcoming community, and we'd love for everyone to bring
their contributions to make it even better.
We appreciate contributions of any size.

* Found a bug or have an idea for a new feature? - [Open an issue](https://github.com/Originate/git-town/issues/new)
* Fixed a bug or created a new feature that others will enjoy? - [Create a pull request](https://help.github.com/articles/using-pull-requests/)

This guide will help you get started and outline some things you should know when developing Git Town.


## Requirements

* Ruby 2.2
  (install [directly](https://www.ruby-lang.org/en/documentation/installation),
  or via a ruby manager like [rvm](https://rvm.io/)
  or [rbenv](https://github.com/sstephenson/rbenv))
  * language the tests are written in
* [ShellCheck](https://github.com/koalaman/shellcheck)
  * used in the linting process to find common errors in the Bash code


## Setup

* install the [requirements](#requirements)
* fork and clone the repository to your machine
* run `bundle` to install ruby gems
* optionally run `rake` to make sure all tests pass on your machine


## Testing

* tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).
* all features need to have comprehensive test coverage
* source code and test files must pass the linters

```bash
# running the different test types
rake         # runs all tests
rake lint    # runs the linters
rake test    # runs the feature tests

# running individual scenarios/features
cucumber <filename>[:<lineno>]
cucumber -n '<scenario/feature name>'

# running individual scenarios/features while showing the application output
DEBUG_COMMANDS=true cucumber <filename>[:<lineno>]

# running several features in parallel
bin/cuke [cucumber parameters]

# auto-fixing formatting issues
rake format  # Run formatters (fixes some lint errors)
```

The `rake [parameters]` commands above can also be run as `bundle exec rake [parameters]`
if you encounter issues.

Git Town's [CI server](https://circleci.com/gh/Originate/git-town)
automatically tests all commits and pull requests,
and notifies you via email and through status badges in pull requests
about problems.


## Pull Requests

Each pull request (PR) should have the same (optional) description that it will have
when committed later.
Besides a general brief description, this also includes
the [issue](https://github.com/Originate/git-town/issues)
it resolves.

When merging approved PRs:
* use `git ship`
* the message for the squashed commit should follow the
  [formatting guidelines for commit messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
* [mention the closed issue](https://help.github.com/articles/closing-issues-via-commit-messages)
in the commit body, so that the respective issue is automatically closed.

Example for a commit message:

```
Automatically foo commits in "git bar"

Fooing changes before they get barred into a separate branch
keeps the final foo simpler when shipping that branch later.

Implements #123
```


## Architecture

*The following refers to all commands except `git-pr`, `git-repo`, and `git-town`.*

Each Git Town command begins by inspecting the current state of the Git repository
(which branch you are on, whether you have open changes).
If there are no errors, it generates a list of steps to run.
Each step is a bash function that wraps an individual Git command.
This list is then executed one by one.

For discussion around this architecture see
[#199](https://github.com/Originate/git-town/issues/199),
where it was proposed.


## Documentation

Every Git Town command
* has a [man page](./man/man1)
* has a [Markdown page](./documentation/commands) that is identical to the man page
* is listed on the [git-town man page](./man/man1/git-town.1)
* is listed on the [README](./README.md)
