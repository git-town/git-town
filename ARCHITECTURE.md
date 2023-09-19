# Git Town architecture

The `cmd` package defines all Git Town commands. Each Git Town command begins by
inspecting the current state of the Git repository (which branches exist,
whether you have open changes). It generates the list of steps and executes them
via Git Town's `runvm`.

The `runstate` package provides facilities to gradually build and represent
steps to execute.

The `steps` package defines all available steps that Git Town can execute.
Examples are steps to change to a different Git branch or to pull updates for
the current branch.

The `runvm` package provides the virtual machine to execute a list of steps.
When executing a step, the runvm.Execute function executes each step in the list
of steps one by one. If a step fails (for example due to a merge conflict), the
engine asks the step to create it's corresponding abort and continue steps, adds
them to the respective StepLists, saves the entire runstate to disk, informs the
user, and exits.

The `persistence` package persists the runstate to disk.

When running "git town continue", Git Town loads the runstate and executes the
"continue" StepList in it. When running "git town abort", Git Town loads the
runstate and executes the "abort" StepList in it. When running "git town undo",
Git Town loads the runstate and executes the "undo" StepList in it.
