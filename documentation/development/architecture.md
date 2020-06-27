# Architecture of the Git Town code base

This mono-repo contains all Git Town related code. The [src](../../src) folder
contains these packages:

- [src/browsers](../../src/browsers) interacts with the local browser
- [src/cli](../../src/cli) reads and writes data from and to Git Town's CLI
- [src/cmd](../../src/cmd) defines Git Town's
  [Cobra](https://github.com/spf13/cobra) subcommands
- [src/config](../../src/config) accesses the Git Town configuration
- [src/drivers](../../src/drivers) interacts with external Git hosting providers
- [src/git](../../src/git) accesses the Git binary on the user's computer
- [src/prompt](../../src/prompt) implements interactive wizards
- [src/run](../../src/run) runs commands in subshells
- [src/steps](../../src/steps) contains the building blocks of Git Town commands
- [src/stringslice](../../src/stringslice) functions for working with string
  slices

## State files

When a Git Town command finishes, it stores the steps to undo or continue itself
in a _state file_ located in the temp folder of the machine. The struct that
holds the state is in [src/steps/run_state.go](../../src/steps/run_state.go).

## Tests

See [test-architecture.md](test-architecture.md).
