# Architecture of the Git Town code base

## Source code

Git Town is written in [Go](https://golang.org).
The source code is located in the [src](../../src) folder.


### CLI wrapper

Git Town uses the [Cobra](https://github.com/spf13/cobra) CLI framework.
[src/cmd](../../src/cmd) contains the commands that Git Town understands.


### Steps

The individual steps that each Git Town command executes are expressed via the
[command pattern](https://en.wikipedia.org/wiki/Command_pattern).
This allows allow for robust and fully automated `--continue`, `--abort`,
and `--undo` functionality.
To distinguish the command-pattern commands from the Git Town commands in [src/cmd](../../src/cmd),
we'll call the former `steps` from now on.
Steps are in [src/steps](../../src/steps)
and implement the individual steps that each Git Town command performs,
like for example [changing to a different Git branch](../../src/steps/checkout_branch_step.go)
or [pulling down updates for the current branch](../../src/steps/pull_branch_step.go).

Each of the possible steps is a Go struct that has a `Run` method to execute the step
as well as a `CreateUndoStepBeforeRun` method
that returns a step that performs the inverse operation.


### Utility code

The other folders in [src](../../src) are utility methods used by the commands and steps:
* [src/browsers](../../src/browsers) provides code to open a browser window with a given URL
* [src/cfmt](../../src/cfmt) contains helpers to print colored text in the terminal
* [src/command](../../src/command) contains a helper to run external tools like Git in a subshell and capture their output
* [src/drivers](../../src/drivers) contains the [driver infrastructure](drivers.md)
  for the APIs of various code hosting services that Git Town supports
* [src/dryrun](../../src/dryrun) contains code that allows to run commands
  that only print but don't execute their steps
* [src/git](../../src/git) contains code to run various Git commands and parse their output intelligently
* [src/prompt](../../src/prompt) contains the code to interactively ask the user for information via the command line
* [src/script](.../../src/script) contains high-level, Git-Town specific helpers
* [src/util](../../src/util) contains a variety of other low-level helper methods


## State files

When a Git Town command finishes, a list of the steps to undo the command
as well as optionally any remaining steps to be run (in case the command aborted)
are stored in state files in a temp folder on the machine.
The code for that is in [src/steps/save_state.go](../../src/steps/save_state.go)
and [src/steps/load_state.go](../../src/steps/load_state.go).


## Tests

The test architecture is described [separately](testing.md).
