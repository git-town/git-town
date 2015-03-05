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
  * Language the tests are written in
* [ShellCheck](https://github.com/koalaman/shellcheck)
  * Used in the linting process to find common errors in the Bash code


## Setup

* fork and clone the repository to your machine
* `bundle` to install ruby gems


## Testing

* Tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).
* All features need to have comprehensive test coverage
* Source code and test files must pass the linters

```bash
# rake tasks
rake         # Run linters and feature tests
rake format  # Run formatters (fixes some lint errors)
rake lint    # Run linters
rake test    # Run feature tests

# run single scenario/feature
cucumber -n 'scenario/feature name'
cucumber [filename][:lineno]

# run single scenario/feature while showing the application output
DEBUG_COMMANDS=true cucumber [filename][:lineno]

# run features in parallel
bin/cuke [<folder>...]
```


## Architecture

*The following refers to all commands except `git-pull-request`, `git-repo`, and `git-town`.*

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
