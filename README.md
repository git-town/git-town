## _Git Automation for Agile Development Teams_
<a href="https://travis-ci.org/Originate/git-town" alt="Build Status" target="_blank"><img src="https://travis-ci.org/Originate/git-town.svg?branch=master"></a>

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

* [git extract](/documentation/git-extract.md) - copy selected commits from the current branch into their own branch
* [git hack](/documentation/git-hack.md) - cut a new feature branch off the main branch
* [git kill](/documentation/git-kill.md) - remove an obsolete feature branch
* [git pr](/documentation/git-pr.md) - create a new pull request
* [git prune-branches](/documentation/git-prune-branches.md) - delete merged branches
* [git repo](/documentation/git-repo.md) - view the repository homepage
* [git ship](/documentation/git-ship.md) - deliver a completed feature branch
* [git sync](/documentation/git-sync.md) - update the current branch with all relevant changes
* [git sync-fork](/documentation/git-sync-fork.md) - pull upstream updates into a forked repository
* [git town](/documentation/git-town.md) - view and change Git Town configuration and easily view help page


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

To install the Fish autocompletion, copy
[Git Town's Fish shell autocompletion definitions file](autocomplete/git.fish)
into your `~/.config/fish/completions/` folder.
Create the folder if necessary.


## Configuration

On first use, Git Town will ask for the main branch name and the names of any other non-feature branches.
Git Town stores its configuration in the Git configuration of your project.
If these ever need to change, the configuration can be updated using [git config](http://git-scm.com/docs/git-config).


```bash
# Read configuration
git config git-town.main-branch-name
git config git-town.non-feature-branch-names

# Write configuration
git config git-town.main-branch-name master
git config git-town.non-feature-branch-names 'qa, production'
```

## Documentation

In addition to the online documentation here,
you can run `git help town` on the command line
for an overview of the git town commands,
or `git help <command>` (e.g. `git help sync`)
for help on an individual command.


## Development

Tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).

```bash
# install tools
bundle
brew install shellcheck  # bash linter

# rake tasks
rake            # Run all linters and specs
rake lint       # Run all linters
rake lint:bash  # Run bash linter
rake lint:ruby  # Run ruby linter
rake spec       # Run specs

# run single test
cucumber -n 'scenario/feature name'
cucumber [filename]:[lineno]

# run cucumber in parallel
bin/cuke [<folder>...]
```

Found a bug or want to contribute a feature?
[Open an issue](https://github.com/Originate/git-town/issues/new)
or - even better - get down, go to town, and fire a feature-tested and linter-passing
[pull request](https://help.github.com/articles/using-pull-requests/)
our way!


## Roadmap

The future roadmap is planned using [GitHub issues](https://github.com/Originate/git-town/issues).
The past roadmap is in the [release notes](release-notes.md).

If you have an idea about a cool feature you would like to see in Git Town,
please [open a ticket](https://github.com/Originate/git-town/issues/new).
Our team will add the [idea](https://github.com/Originate/git-town/labels/idea) tag.
Once we reach agreement about this idea, it will be tagged as an
[enhancement](https://github.com/Originate/git-town/labels/enhancement)
or a [bug](https://github.com/Originate/git-town/labels/bug).
