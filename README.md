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

#### Typical development commands

- [git hack](https://www.git-town.com/commands/hack.html) - create a new
  up-to-date feature branch off the main branch
- [git sync](https://www.git-town.com/commands/sync.html) - update existing
  branches, remove shipped branches
- [git propose](https://www.git-town.com/commands/propose.html) - create a pull
  or merge request for a feature branch

#### Stacked changes

- [git append](https://www.git-town.com/commands/append.html) - insert a new
  branch as a child of the current branch
- [git prepend](https://www.git-town.com/commands/prepend.html) - insert a new
  branch between the current branch and its parent
- [git set-parent](https://www.git-town.com/commands/set-parent.html) - update
  the parent of a branch

#### Branch Types

- [git observe](https://www.git-town.com/commands/observe.html) - track somebody
  else's feature branch without contributing to it
- [git contribute](https://www.git-town.com/commands/observe.html) - add commits
  to somebody else's feature branch
- [git park](https://www.git-town.com/advanced-syncing#parked-branches) - stop
  syncing one of your feature branches

#### Advanced development commands

- [git kill](https://www.git-town.com/commands/kill.html) - remove a feature
  branch
- [git rename-branch](https://www.git-town.com/commands/rename-branch.html) -
  rename a branch
- [git repo](https://www.git-town.com/commands/repo.html) - view the repository
  homepage
- [git ship](https://www.git-town.com/commands/ship.html) - merge a completed
  feature branch and remove it

#### Configuration

- [git town config](https://www.git-town.com/commands/config.html) - display or
  update your Git Town configuration
- [git town version](https://www.git-town.com/commands/version.html) - display
  the installed version of Git Town

## Installation

See the [installation](https://www.git-town.com/install.html) and
[configuration](https://www.git-town.com/quick-configuration.html) instructions.

## Documentation

The [Git Town website](https://www.git-town.com) provides documentation for Git
Town users. `git town help [command]` shows help on the CLI.

## Contributing

Found a bug or have an idea for a new feature?
[Open an issue](https://github.com/git-town/git-town/issues/new) or send a
[pull request](https://help.github.com/articles/using-pull-requests)! Our
[developer documentation](docs/DEVELOPMENT.md) helps you get started.

[![Stargazers over time](https://starchart.cc/git-town/git-town.svg)](https://starchart.cc/git-town/git-town)
