## _Git Automation for Agile Development Teams_
<a href="https://travis-ci.org/Originate/git-town" alt="Build Status" target="_blank"><img src="https://travis-ci.org/Originate/git-town.svg?branch=master"></a>
[![License](http://img.shields.io/:license-MIT-blue.svg?style=flat)](MIT-LICENSE)

Git Town provides a number of additional Git commands that
automate the typical high-level operations in
[GitHub Flow](http://scottchacon.com/2011/08/31/github-flow.html)
and other workflows.

It is designed for workflows that have a main branch
(typically "development" or "master")
from which feature branches are cut and into which they are merged,
and it assumes you use a central code repository like
[GitHub](http://github.com/) or [Bitbucket](https://bitbucket.org/).


## Commands

* [git extract](/documentation/commands/git-extract.md) - copy selected commits from the current branch into their own branch
* [git hack](/documentation/commands/git-hack.md) - cut a new feature branch off the main branch
* [git kill](/documentation/commands/git-kill.md) - remove an obsolete feature branch
* [git pr](/documentation/commands/git-pr.md) - create a new pull request
* [git prune-branches](/documentation/commands/git-prune-branches.md) - delete merged branches
* [git repo](/documentation/commands/git-repo.md) - view the repository homepage
* [git ship](/documentation/commands/git-ship.md) - deliver a completed feature branch
* [git sync](/documentation/commands/git-sync.md) - update the current branch with all relevant changes
* [git sync-fork](/documentation/commands/git-sync-fork.md) - pull upstream updates into a forked repository
* [git town](/documentation/commands/git-town.md) - general Git Town help, view and change Git Town configuration


#### Notes

* minimizes network requests
  * each command performs a single fetch
  * skips unnecessary pushes
* automatically prunes deleted remote branches


## Installation

Git Town is written in Bash, so it runs anywhere Git and Bash runs.
Installation on OS X can be done using [Homebrew](http://brew.sh/).
Other platforms need to install manually.

<table>
  <tr>
    <th width="300px">
      Using Homebrew
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
      <ol>
        <li>clone the repo to your machine (into DIR)</li>
        <li>add DIR/src to your <code>$PATH</code></li>
        <li>add DIR/man to your <code>$MANPATH</code></li>
        <li>
          install <a href="http://en.wikipedia.org/wiki/Dialog_(software)">Dialog</a>
          (used by <a href="/documentation/git-extract.md">git extract</a>)
        </li>
      </ol>
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
      <ol>
        <li>run <code>git pull</code> in DIR</li>
      </ol>
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
      <ol>
        <li>remove DIR</li>
        <li>remove DIR/src from your <code>$PATH</code></li>
        <li>remove DIR/man from your <code>$MANPATH</code></li>
      </ol>
    </td>
  </tr>
</table>


#### Optional tools that make Git Town better

* __Autocompletion for [Fish shell](http://fishshell.com)__

    ```
    mkdir -p ~/.config/fish/completions/
    curl -o ~/.config/fish/completions/git.fish http://raw.githubusercontent.com/Originate/git-town/master/autocomplete/git.fish
    ```


## Configuration

Git Town is configured on a per-repository basis. Upon first use in a given repository, Git Town will ask the user for all required
configuration information. Use the [`git town`](/documentation/git-town.md) command to view and update your configuration at any time.


## Documentation

In addition to the online documentation here,
you can run `git help town` on the command line
for an overview of the git town commands,
or `git help <command>` (e.g. `git help sync`)
for help on an individual command.


## Contributing

Found a bug or have an idea for a new feature?
[Open an issue](https://github.com/Originate/git-town/issues/new)
or - even better - get down, go to town, and fire a feature-tested and linter-passing
[pull request](https://help.github.com/articles/using-pull-requests/)
our way!

Check out our [development notes](/documentation/development.md) to start coding.
