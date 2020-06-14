# Architecture of the Git Town code base

All Git Town code is stored as a mono-repo. The [Go](https://golang.org) source
code is located in the [src](../../src) folder. It contains these packages:

- [src/browsers](../../src/browsers) interacts with the local browser
- [src/cli](../../src/cli) reads and writes data from and to Git Town's CLI
- [src/cmd](../../src/cmd) defines Git Town's
  [Cobra](https://github.com/spf13/cobra)-based subcommands
- [src/command](../../src/command) runs commands in subshells
- [src/config](../../src/config) accesses the Git Town configuration
- [src/drivers](../../src/drivers) interacts with external Git hosting providers
- [src/dryrun](../../src/dryrun) manages dry-runs
- [src/git](../../src/git) accesses the Git binary on the user's computer
- [src/prompt](../../src/prompt) implements interactive wizards
- [src/steps](../../src/steps) contains the building blocks of Git Town commands
- [src/util](../../src/util) string helper functions

## State files

When a Git Town command finishes, a list of the steps to undo the command as
well as optionally any remaining steps to be run (in case the command
encountered conflicts) are stored in a state files in a temp folder on the
machine. The struct that holds the state is in
[src/steps/run_state.go](../../src/steps/run_state.go).

## Tests

The test architecture is described [separately](testing.md).
