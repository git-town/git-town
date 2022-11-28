<p align="center">
  <img src="https://raw.githubusercontent.com/git-town/git-town/main/website/src/logo.svg">
  <br>
  <img src="https://github.com/git-town/git-town/actions/workflows/cuke.yml/badge.svg" alt="end-to-end test status">
  <img src="https://github.com/git-town/git-town/actions/workflows/lint_unit.yml/badge.svg" alt="unit test status">
  <a href="https://goreportcard.com/report/github.com/git-town/git-town">
    <img src="https://goreportcard.com/badge/github.com/git-town/git-town" alt="Go report card status">
  </a>
  <a href="https://www.codetriage.com/originate/git-town">
    <img src="https://www.codetriage.com/originate/git-town/badges/users.svg" alt="Help Contribute to Open Source">
  </a>
  <img src="https://api.netlify.com/api/v1/badges/c2ea5505-be48-42e5-bb8a-b807d18d99ed/deploy-status" alt="Netlify deploy status">
</p>

Git Town makes [Git](https://git-scm.com) more efficient, especially for large
teams. See [this screencast](https://youtu.be/oLaUsUlFfTo) for an introduction.

## Commands

Git Town provides these additional Git commands:

#### Typical development commands

- [git hack](https://www.git-town.com/commands/hack.html) - cuts a new
  up-to-date feature branch off the main branch
- [git sync](https://www.git-town.com/commands/sync.html) - updates the current
  branch with all ongoing changes
- [git new-pull-request](https://www.git-town.com/commands/new-pull-request.html) -
  create a new pull request
- [git ship](https://www.git-town.com/commands/ship.html) - delivers a completed
  feature branch and removes it

#### Advanced development commands

- [git kill](https://www.git-town.com/commands/kill.html) - removes a feature
  branch
- [git prune-branches](https://www.git-town.com/commands/prune-branches.html) -
  delete all merged branches
- [git rename-branch](https://www.git-town.com/commands/rename-branch.html) -
  rename a branch
- [git repo](https://www.git-town.com/commands/repo.html) - view the repository
  homepage

#### Nested feature branches

- [git append](https://www.git-town.com/commands/append.html) - insert a new
  branch as a child of the current branch
- [git prepend](https://www.git-town.com/commands/prepend.html) - insert a new
  branch between the current branch and its parent
- [git set-parent-branch](https://www.git-town.com/commands/set-parent-branch.html) -
  updates a branch's parent

#### Git Town configuration

- [git town config](https://www.git-town.com/commands/config.html) - displays or
  updates your Git Town configuration
- [git town new-branch-push-flag](https://www.git-town.com/commands/new-branch-push-flag.html) -
  configures whether new empty branches get pushed to origin
- [git town main-branch](https://www.git-town.com/commands/main-branch.html) -
  displays or sets the main development branch for the current repo
- [git town offline](https://www.git-town.com/commands/offline.html) -
  enables/disables offline mode
- [git town perennial-branches](https://www.git-town.com/commands/perennial-branches.html) -
  displays or updates the perennial branches for the current repo
- [git town pull-branch-strategy](https://www.git-town.com/commands/pull-branch-strategy.html) -
  displays or sets the strategy to update perennial branches

#### Git Town setup

- [git town alias](https://www.git-town.com/commands/alias.html) - adds or
  removes shorter aliases for Git Town commands
- [git town completions](https://www.git-town.com/commands/completions.html) -
  generates completion scripts for Bash, zsh, fish & PowerShell.
- [git town version](https://www.git-town.com/commands/version.html) - displays
  the installed version of Git Town

## Installation

See the [installation](https://www.git-town.com/install.html) and
[configuration](https://www.git-town.com/quick-configuration.html) instructions
for more details.

## Documentation

The [Git Town website](https://www.git-town.com) provides documentation for Git
Town users. `git town help [command]` shows help on the CLI.

## Contributing

Found a bug or have an idea for a new feature?
[Open an issue](https://github.com/git-town/git-town/issues/new) or send a
[pull request](https://help.github.com/articles/using-pull-requests)! Our
[developer documentation](DEVELOPMENT.md) helps you get started.

[![Stargazers over time](https://starchart.cc/git-town/git-town.svg)](https://starchart.cc/git-town/git-town)
