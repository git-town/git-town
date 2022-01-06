# Commands

Git Town adds the following commands to Git. Commands are provided in the short
form, similar to how `git town alias` creates them.

### Development Workflow

- [git hack](/documentation/commands/hack.md) - create a new feature branch
- [git sync](/documentation/commands/sync.md) - update the current branch with
  all ongoing changes
- [git new-pull-request](/documentation/commands/new-pull-request.md) - create a
  new pull request
- [git ship](/documentation/commands/ship.md) - deliver a completed feature
  branch

### Repository Maintenance

- [git kill](/documentation/commands/kill.md) - delete a feature branch
- [git prune-branches](/documentation/commands/prune-branches.md) - remove all
  merged branches
- [git rename-branch](/documentation/commands/rename-branch.md) - rename a
  branch
- [git append](/documentation/commands/append.md) - create a new feature branch
  as a child of the current branch
- [git prepend](/documentation/commands/prepend.md) - create a new feature
  branch between the current branch and its parent
- [git repo](/documentation/commands/repo.md) - view the Git repository in the
  browser

### Git Town Configuration

- [git town config](/documentation/commands/config.md) - display or update your
  Git Town configuration
- [git town new-branch-push-flag](/documentation/commands/new-branch-push-flag.md) -
  configure whether new empty branches are pushed to origin
- [git town main-branch](/documentation/commands/main-branch.md) - display/set
  the main development branch for the current repo
- [git town offline](/documentation/commands/offline.md) - enable/disable
  offline mode
- [git town perennial-branches](/documentation/commands/perennial-branches.md) -
  display or update the perennial branches for the current repo
- [git town pull-branch-strategy](/documentation/commands/pull-branch-strategy.md) -
  display or set the strategy with which perennial branches are updated
- [git town set-parent-branch](/documentation/commands/set-parent-branch.md) -
  change the parent of a feature branch

### Git Town Installation

- [git town alias](/documentation/commands/alias.md) - add or remove shorter
  aliases for Git Town commands
- [git town completions](/documentation/commands/completions.md) - generate
  completion scripts for Bash, zsh, fish & PowerShell.
- [git town version](/documentation/commands/version.md) - display the installed
  version of Git Town
