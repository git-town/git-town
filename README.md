![Git Town](http://originate.github.io/git-town/documentation/logo-horizontal.svg)

[![Build Status](https://circleci.com/gh/Originate/git-town/tree/master.svg?style=shield)](https://circleci.com/gh/Originate/git-town/tree/master)
[![License](http://img.shields.io/:license-MIT-blue.svg?style=flat)](LICENSE)
[![Join the chat at https://gitter.im/Originate/git-town](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/Originate/git-town?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Git Town makes software development teams who use Git even more productive and happy.
It adds additional Git commands that support
[GitHub Flow](http://scottchacon.com/2011/08/31/github-flow.html),
[Git Flow](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow),
the [Nvie model](http://nvie.com/posts/a-successful-git-branching-model),
[GitLab Flow](https://about.gitlab.com/2014/09/29/gitlab-flow/),
and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.

See <http://www.git-town.com> for documentation.


## Commands

Git Town provides these additional Git commands:

__Development Workflow__

* [git hack](/documentation/commands/git-hack.md) - cuts a new up-to-date feature branch off the main branch
* [git sync](/documentation/commands/git-sync.md) - updates the current branch with all ongoing changes
* [git new-pull-request](/documentation/commands/git-new-pull-request.md) - create a new pull request
* [git ship](/documentation/commands/git-ship.md) - delivers a completed feature branch and removes it


__Repository Maintenance__

* [git extract](/documentation/commands/git-extract.md) - copies selected commits from the current branch into their own branch
* [git kill](/documentation/commands/git-kill.md) - removes a feature branch
* [git prune-branches](/documentation/commands/git-prune-branches.md) - delete all merged branches
* [git rename-branch](/documentation/commands/git-rename-branch.md) - rename a branch
* [git repo](/documentation/commands/git-repo.md) - view the repository homepage
* [git sync-fork](/documentation/commands/git-sync-fork.md) - pull upstream updates into a forked repository


__Configuration and Help__

* [git town](/documentation/commands/git-town.md) - general Git Town help, view and change Git Town configuration


## Installation

Git Town runs anywhere Git and [Bash](https://www.gnu.org/software/bash/bash.html) run.
Check out our [installation instructions](documentation/installation.md) for more details.


## Configuration

Git Town is configured on a per-repository basis. It requires two pieces of information.

* the main development branch
* the perennial branches (see more [here](/documentation/development/branch_hierarchy.md#perennial-branches))

Upon first use in a repository, you will be prompted for this information.
Use the [git town](/documentation/commands/git-town.md) command to view or update your configuration at any time.


## Documentation

In addition to the online documentation here,
you can run `git town` on the command line for an overview of the Git Town commands,
or `git help <command>` (e.g. `git help sync`) for help with an individual command.


## Contributing

Found a bug or have an idea for a new feature?
[Open an issue](https://github.com/Originate/git-town/issues/new)
or - even better - get down, go to town, and fire a feature-tested
[pull request](https://help.github.com/articles/using-pull-requests/)
our way! Check out our [contributing guide](/CONTRIBUTING.md) to start coding.
