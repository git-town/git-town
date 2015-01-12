# Git Town Development

## Requirements

* Ruby 2.1.5
  (install [directly](https://www.ruby-lang.org/en/documentation/installation),
  or via a ruby manager like [rvm](https://rvm.io/)
  or [rbenv](https://github.com/sstephenson/rbenv))
* [ShellCheck](https://github.com/koalaman/shellcheck)


## Setup

* fork and clone the repository to your machine
* `bundle` to install ruby gems


## Running Tests

* Tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).
* All features have need to have comprehensive test coverage
* We are using linters for both the source and test files

```bash
# rake tasks
rake          # Run linters and specs
rake format   # Run formatters (fixes some lint errors)
rake lint     # Run linters
rake spec     # Run specs

# run single test
cucumber -n 'scenario/feature name'
cucumber [filename]:[lineno]

# run cucumber in parallel
bin/cuke [<folder>...]
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


## Command documentation

Every Git Town command
* has a [man page](../man/man1)
* has a [markdown page](./commands) that is identical to the man page
* is listed on the [git-town man page](../man/man1/git-town.1)
* is listed on the [README](../README.md)
