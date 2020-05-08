# Step Lists

_The following refers to all commands except `git-new-pull-request`, `git-repo`
and `git-town`._

The individual steps that each Git Town command executes are represested via the
[command pattern](https://en.wikipedia.org/wiki/Command_pattern). This allows
fully automated and robust `--continue`, `--abort`, and `--undo` functionality
for each command.

To distinguish the command-pattern commands from the Git Town commands in
[src/cmd](../../src/cmd), we'll call the former `steps` from now on.

## Running commands

Each Git Town command begins by inspecting the current state of the Git
repository (which branch you are on, whether you have open changes). If there
are no errors, it generates a list of steps to run.

Steps, located in [src/steps](../../src/steps), implement the individual steps
that each Git Town command performs, like for example
[changing to a different Git branch](../../src/steps/checkout_branch_step.go) or
[pulling down updates for the current branch](../../src/steps/pull_branch_step.go).
They are Go structs that have a `Run` method which executes the step.

When executing a step, the undo steps for it are determined and appended to the
`undo list` for the current Git Town command. This is done by calling the
methods `step.CreateUndoStep()`.

If a Git command fails (typically due to a merge conflict), then the program
saves state (lists with steps to abort and continue) to disk, informs the user
how to abort/continue the current step, and exits.

#### Continuing commands

The step to continue is determined by calling the `CreateContinueStep` method of
the current (failed) step. On the user resolving the issue and continuing, Git
Town runs the continue step and the remaining step list.

`git town sync` also allows the user to skip the current branch, which skips all
commands until the next checkout and then resumes executing steps.

#### Aborting commands

The step to abort is determined by calling the `CreateAbortStep` method of the
current (failed) step. On abort, Git Town executes the abort step and the list
of undo steps for all previously run steps.

#### Undoing commands

A successfully finished command can be undone. To do that, Git Town executes all
the undo steps.
