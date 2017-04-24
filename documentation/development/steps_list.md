# Step Lists

_The following refers to all commands except `git-new-pull-request`, `git-repo` and `git-town`._

Each Git Town command begins by inspecting the current state of the Git repository
(which branch you are on, whether you have open changes).
If there are no errors, it generates a list of steps to run.
Each step is a go `struct`.

When executing a step, the undo steps are determined and added to a separate list.
This is done by calling the methods `step.CreateUndoStepBeforeRun()` and `step.CreateUndoStepAfterRun()`.


## Abort / Continue

If a Git command fails (typically due to a merge conflict), then the program halts
and asks the user what they would like to do. In most cases they can either abort or continue.

If the user aborts then `step.CreateAbortStep()` is called and the undo steps are executed.

If the user resolves the issue and continues then `step.CreateContinueStep()` is called
and we resume executing steps.

`git town-sync` also allows the user to skip the current branch,
which skips all commands until the next checkout and then resumes executing steps.


## Undo

If a command finished successfully, then it can be undone.
This will simply execute all the undo steps.
