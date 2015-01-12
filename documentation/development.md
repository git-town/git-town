# Git Town Development

## Requirements

* Ruby 2.1.5
  (install [directly](https://www.ruby-lang.org/en/documentation/installation),
  or via a ruby manager like [rvm](https://rvm.io/)
  or [rbenv](https://github.com/sstephenson/rbenv))
* [ShellCheck](https://github.com/koalaman/shellcheck)


## Setup

* fork and clone repository to your machine
* `bundle` to install ruby gems


## Running Tests

* Tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).
* We have linters for bash (source), ruby (tests), and cucumber (tests)

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

Git Town commands are simply running a series of git commands.
Each command inspects the current state of the git repository
(which branch you are on, do you have open changes)
and generates a list of steps to run.
Each step is a bash function that is a wrapper around a git command.
This list is then executed one by one.

For discussion around this architecture see
[#199](https://github.com/Originate/git-town/issues/199)
where it was proposed.


## Command documentation

Each command has:
* [a man page](../man/man1)
* [a markdown page](./commands)
* is listed on the [git-town man page](../man/man1/git-town.1)
* is listed on the [README](../README.md)
