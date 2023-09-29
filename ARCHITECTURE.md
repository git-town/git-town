# Git Town architecture

The [cmd](src/cmd) package defines all Git Town commands. Each Git Town command
inspects the current state of the Git repository (branches, Git configuration,
stash size) and generates a snapshot of the repository's initial status. It
generates a [list](src/runstate) of [steps](src/steps) that make up the
functionality of the Git Town command and executes them via Git Town's
[virtual machine](src/runvm). Steps call Git via the [git](src/git) package.
When the virtual machine is done, it generates the final status of the Git
repository, compares it to the initial snapshot, and renders the diff into steps
to undo the changes that the Git Town command made.

The `git town continue` command continues executing the
[persisted](src/persistence) runstate. The `git town undo` command executes the
persisted undo step list.
