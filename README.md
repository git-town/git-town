![Git Town](http://originate.github.io/git-town/documentation/logo-horizontal.svg)

[![Build Status](https://circleci.com/gh/Originate/git-town/tree/master.svg?style=shield)](https://circleci.com/gh/Originate/git-town/tree/master)
[![License](http://img.shields.io/:license-MIT-blue.svg?style=flat)](LICENSE)

Git Town makes software development teams who use Git even more productive and happy.
It adds additional Git commands that support
[GitHub Flow](http://scottchacon.com/2011/08/31/github-flow.html),
[Git Flow](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow),
the [Nvie model](http://nvie.com/posts/a-successful-git-branching-model),
[GitLab Flow](https://about.gitlab.com/2014/09/29/gitlab-flow/),
and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.

Check out [the big picture](documentation/background.md) for more background on Git Town
and the [tutorial](documentation/tutorial.md) to get an idea for how it works.


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

Git Town is written in 100% [Bash](https://www.gnu.org/software/bash/bash.html),
so it runs anywhere Git and Bash run.

<table>
  <tr>
    <th width="300px">
      Using <a href="http://brew.sh">Homebrew</a>
    </th>
    <th width="400px">
      Manually
    </th>
  </tr>
  <tr>
    <td colspan="2">
      <b>Install</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew tap Originate/gittown</code><br>
      <code>brew install git-town</code>
    </td>
    <td>
      <ul>
        <li>clone the repo to your machine (into DIR)</li>
        <li>add DIR/src to your <code>$PATH</code></li>
        <li>add DIR/man to your <code>$MANPATH</code></li>
        <li>
          install <a href="http://en.wikipedia.org/wiki/Dialog_(software)">Dialog</a>
          (used by <a href="/documentation/commands/git-extract.md">git extract</a>)
        </li>
      </ul>
    </td>
  </tr>
  <tr>
    <td colspan="2">
      <b>Update</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew update</code><br>
      <code>brew upgrade git-town</code>
    </td>
    <td>
      <ul>
        <li>run <code>git pull</code> in DIR</li>
      </ul>
    </td>
  </tr>
  <tr>
    <td colspan="2">
      <b>Uninstall</b>
    </td>
  </tr>
  <tr>
    <td>
      <code>brew uninstall git-town</code><br>
      <code>brew untap Originate/gittown</code>
    </td>
    <td>
      <ul>
        <li>remove DIR</li>
        <li>remove DIR/src from your <code>$PATH</code></li>
        <li>remove DIR/man from your <code>$MANPATH</code></li>
      </ul>
    </td>
  </tr>
</table>


#### Optional tools that make Git Town better

* __Autocompletion for [Fish shell](http://fishshell.com)__

    ```
    $ git town install-fish-autocompletion
    ```


## Configuration

Git Town is configured on a per-repository basis.
Upon first use in a repository, it will ask for all required configuration.
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
