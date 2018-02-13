# Step Lists

_The following refers to all commands except `git-new-pull-request`, `git-repo` and `git-town`._

The individual steps that each Git Town command executes are represested via the
[command pattern](https://en.wikipedia.org/wiki/Command_pattern).
This allows allow for robust and fully automated `--continue`, `--abort`,
and `--undo` functionality.
To distinguish the command-pattern commands from the Git Town commands in [src/cmd](../../src/cmd),
we'll call the former `steps` from now on.

Each Git Town command begins by inspecting the current state of the Git repository
(which branch you are on, whether you have open changes).
If there are no errors, it generates a list of steps to run.
Steps are in [src/steps](../../src/steps)
and implement the individual steps that each Git Town command performs,
like for example [changing to a different Git branch](../../src/steps/checkout_branch_step.go)
or [pulling down updates for the current branch](../../src/steps/pull_branch_step.go).

Steps are Go structs that have a `Run` method to execute the step
as well as a `CreateUndoStepBeforeRun` method
that returns a step that performs the inverse operation.
When executing a step, the undo steps are determined and added to a separate list.
This is done by calling the methods `step.CreateUndoStepBeforeRun()` and `step.CreateUndoStepAfterRun()`.

## Abort / Continue

If a Git command fails (typically due to a merge conflict), then the program halts
and asks the user what they would like to do. In most cases they can either abort or continue.

If the user aborts then `step.CreateAbortStep()` is called and the undo steps are executed.

If the user resolves the issue and continues then `step.CreateContinueStep()` is called
and we resume executing steps.

`git town sync` also allows the user to skip the current branch,
which skips all commands until the next checkout and then resumes executing steps.

## Undo

If a command finished successfully, then it can be undone.
This will simply execute all the undo steps.
