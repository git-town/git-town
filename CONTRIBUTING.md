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
cucumber [filename]:[lineno]

# run features in parallel
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


## Documentation

Every Git Town command
* has a [man page](../man/man1)
* has a [Markdown page](./commands) that is identical to the man page
* is listed on the [git-town man page](../man/man1/git-town.1)
* is listed on the [README](../README.md)
