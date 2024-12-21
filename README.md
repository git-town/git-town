<p align="center">
  <picture>
    <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/git-town/git-town/main/website/src/logo.svg">
    <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/git-town/git-town/main/website/src/logo-dark.svg">
    <img alt="Git Town logo" src="https://raw.githubusercontent.com/git-town/git-town/main/website/src/logo.svg">
  </picture>
  <br>
  <img src="https://github.com/git-town/git-town/actions/workflows/cuke.yml/badge.svg" alt="end-to-end test status">
  <img src="https://github.com/git-town/git-town/actions/workflows/unit.yml/badge.svg" alt="unit test status">
  <img src="https://github.com/git-town/git-town/actions/workflows/lint_docs.yml/badge.svg" alt="linters and documentation test status">
  <img src="https://github.com/git-town/git-town/actions/workflows/windows.yml/badge.svg" alt="windows tests">
  <a href="https://goreportcard.com/report/github.com/git-town/git-town"><img src="https://goreportcard.com/badge/github.com/git-town/git-town" alt="Go report card status"></a>
  <img src="https://api.netlify.com/api/v1/badges/c2ea5505-be48-42e5-bb8a-b807d18d99ed/deploy-status" alt="Netlify deploy status">
</p>

Git Town provides additional Git commands that automate the creation,
synchronization, shipping, and cleanup of Git branches. Compatible with all
popular Git workflows like Git Flow, GitHub Flow, GitLab Flow, and trunk-based
development. Supports mono-repos and stacked changes. Check out
[this screencast](https://youtu.be/oLaUsUlFfTo) for an introduction.

#### Basic development commands

- [hack](https://www.git-town.com/commands/hack.html) - create a new up-to-date
  feature branch off the main branch
- [sync](https://www.git-town.com/commands/sync.html) - update existing
  branches, remove shipped branches
- [switch](https://www.git-town.com/commands/switch.html) - switch between
  branches via text UI
- [propose](https://www.git-town.com/commands/propose.html) - create a pull or
  merge request for a feature branch

#### Stacked changes

- [append](https://www.git-town.com/commands/append.html) - insert a new branch
  as a child of the current branch
- [diff-parent](https://www.git-town.com/commands/diff-parent.html) - show the
  changes committed to a feature branch
- [merge](https://www.git-town.com/commands/merge.html) - merge two adjacent
  branches in a stack into one
- [prepend](https://www.git-town.com/commands/prepend.html) - insert a new
  branch between the current branch and its parent
- [set-parent](https://www.git-town.com/commands/set-parent.html) - update the
  parent of a branch

#### Limit branch syncing

- [contribute](https://www.git-town.com/commands/contribute) - add commits to
  somebody else's feature branch
- [observe](https://www.git-town.com/commands/observe) - track somebody else's
  feature branch without contributing to it
- [park](https://www.git-town.com/commands/park) - stop syncing one of your
  feature branches
- [prototype](https://www.git-town.com/commands/prototype) - sync but don't push
  a branch

#### Dealing with errors

- [continue](https://www.git-town.com/commands/continue.html) - resume the last
  run Git Town command after having resolved conflicts
- [skip](https://www.git-town.com/commands/skip.html) - resume the last run Git
  Town command by skipping the current branch
- [status](https://www.git-town.com/commands/status.html) - displays or resets
  the current suspended Git Town command
- [undo](https://www.git-town.com/commands/undo.html) - undo the most recent Git
  Town command

#### Setup and configuration

- [config](https://www.git-town.com/commands/config.html) - display or update
  your Git Town configuration
- [config setup](https://www.git-town.com/commands/config-setup) - run the
  visual setup assistant
- [offline](https://www.git-town.com/commands/offline.html) - start or stop
  running in offline mode
- [completions](https://www.git-town.com/commands/completions) -
  auto-completions for bash, zsh, fish, and PowerShell

#### Advanced development commands

- [branch](https://www.git-town.com/commands/branch) - display the local branch
  hierarchy
- [compress](https://www.git-town.com/commands/compress.html) - squash all
  commits on feature branches down to a single commit
- [delete](https://www.git-town.com/commands/delete.html) - remove a feature
  branch
- [rename](https://www.git-town.com/commands/rename.html) - rename a branch
- [repo](https://www.git-town.com/commands/repo.html) - view the repository
  homepage
- [ship](https://www.git-town.com/commands/ship.html) - merge a completed
  feature branch and remove it

## Installation

See the [installation](https://www.git-town.com/install.html) and
[configuration](https://www.git-town.com/configuration) instructions.

## Documentation

The [Git Town website](https://www.git-town.com) provides documentation for Git
Town users. `git town help [command]` shows help on the CLI.

## Contributing

Found a bug or have an idea for a new feature?
[Open an issue](https://github.com/git-town/git-town/issues/new) or send a
[pull request](https://help.github.com/articles/using-pull-requests)! Our
[developer documentation](docs/DEVELOPMENT.md) helps you get started.

[![Stargazers over time](https://starchart.cc/git-town/git-town.svg)](https://starchart.cc/git-town/git-town)
