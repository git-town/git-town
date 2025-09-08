# Commands

Run `git town` for an overview of all Git Town commands and
`git town help <command>` for help with individual commands.

### Basic workflow

- [git town hack](commands/hack.md) - create a new feature branch
- [git town sync](commands/sync.md) - update the current branch with all ongoing
  changes
- [git town switch](commands/switch.md) - switch between branches visually
- [git town propose](commands/propose.md) - propose to ship a branch

### Dealing with errors

- [git town continue](commands/continue.md) - continue after you resolved the
  merge conflict
- [git town skip](commands/skip.md) - when syncing all branches, ignore the
  current branch and continue with the next one
- [git town status](commands/status.md) - display available commands
- [git town undo](commands/undo.md) - undo the last completed Git Town command

### Stacked changes

- [git town append](commands/append.md) - create a new feature branch as a child
- [git town detach](commands/detach.md) - move a branch out of a stack
- [git town down](commands/down.md) - switch to a child of the current branch
- [git town diff-parent](commands/diff-parent.md) - show the changes committed
  to a branch
- [git town merge](commands/merge.md) - merges the current branch with its
  parent
- [git town prepend](commands/prepend.md) - create a new feature branch between
  the current branch and its parent
- [git town set-parent](commands/set-parent.md) - change the parent of a feature
  branch
- [git town swap](commands/swap.md) - swap the position of this branch with its
  parent
- [git town up](commands/up.md) - switch to the parent of the current stack

### Limit branch syncing

- [git town contribute](commands/detach.md) - stop syncing some feature branches
  with their parents
- [git town observe](commands/observe.md) - stop syncing your contributions to
  some feature branches
- [git town park](commands/park.md) - suspend syncing some of your feature
  branches altogether
- [git town prototype](commands/prototype.md) - sync but don't push a local
  feature branch

### Git Town setup

- [git town completion](commands/completions.md) - generate completion scripts
  for Bash, zsh, fish & PowerShell.
- [git town config](commands/config.md) - display or update your Git Town
  configuration
- [git town init](commands/init.md) - setup assistant
- [git town offline](commands/offline.md) - enable/disable offline mode

### Additional commands

- [git town branch](commands/branch.md) - display the local branch hierarchy and
  types
- [git town compress](commands/compress.md) - squash all commits on a branch
  branches down to a single commit
- [git town delete](commands/delete.md) - delete a feature branch
- [git town help](commands/help.md) - help about any command
- [git town rename](commands/rename.md) - rename a branch
- [git town repo](commands/repo.md) - view the Git repository in the browser
- [git town ship](commands/ship.md) - deliver a completed feature branch
