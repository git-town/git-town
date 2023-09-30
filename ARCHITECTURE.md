# Git Town architecture

The [cmd](src/cmd) package defines all Git Town commands. Each Git Town command
inspects the current state of the Git repository (branches, Git configuration,
stash size) and generates a snapshot of the repository's initial status. It
generates a [list](src/runstate) of [steps](src/steps) that make up the
functionality of the Git Town command and executes them via Git Town's
[virtual machine](src/runvm). When the virtual machine is done, it generates
another snapshot of the Git repository, compares it to the initial snapshot, and
renders the diff into steps to undo the changes that the Git Town command made.

Steps call Git via the [git](src/git) package which in turn uses the
[subshell](src/subshell) package. Communication with source code hosting
providers happens in the [hosting](src/hosting) package. Git Town
remote-controls browsers via the [browser](src/browser) package and queries
additional information from the user via the [dialog](src/dialog) package.

The `git town continue` command continues executing the
[persisted](src/persistence) runstate. The `git town undo` command executes the
persisted undo step list.
