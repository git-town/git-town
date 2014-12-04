## _Git Automation for Agile Development Teams_
<a href="https://travis-ci.org/Originate/git-town" alt="Build Status" target="_blank"><img src="https://travis-ci.org/Originate/git-town.svg?branch=master"></a>

Git Town provides a number of additional Git commands that
automate the typical high-level operations in
<a href="http://scottchacon.com/2011/08/31/github-flow.html" target="_blank">GitHub flow</a>
and others.

It is designed for workflows that have a main branch
(typically "development" or "master")
from which feature branches are cut and into which they are merged,
and assumes you use a central code repository like
<a href="http://github.com" target="_blank">GitHub</a> or
<a href="https://bitbucket.org" target="_blank">BitBucket</a>.


## commands

* [git extract](/documentation/git-extract.md) - copy selected commits from the current branch into their own branch
* [git hack](/documentation/git-hack.md) - cut a new feature branch off the main branch
* [git kill](/documentation/git-kill.md) - remove an obsolete feature branch
* [git pr](/documentation/git-pr.md) - create a new pull request
* [git prune-branches](/documentation/git-prune-branches.md) - delete merged branches
* [git ship](/documentation/git-ship.md) - deliver a completed feature branch
* [git sync](/documentation/git-sync.md) - update the current branch with all relevant changes
* [git sync-fork](/documentation/git-sync-fork.md) - pull upstream updates into a forked repository


#### notes

* minimizes network requests
  * each command performs a single fetch
  * skips unnecessary pushes
* automatically prunes deleted remote branches


## installation

Git Town is 100% Bash script, so it runs anywhere where Git and Bash runs.
Installation on OS X can be done using <a href="http://brew.sh" target="_blank">Homebrew</a>,
other platforms need to install manually.

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
      brew tap Originate/gittown<br>
      brew install git-town
    </td>
    <td>
      clone the repo to your machine (into DIR)<br>
      add DIR to your path<br>
      add DIR/man to your manpath
    </td>
  </tr>
  <tr>
    <td colspan="2">
      <b>Update</b>
    </td>
  </tr>
  <tr>
    <td>
      brew update<br>
      brew upgrade git-town
    </td>
    <td>
      run <code>git pull</code> in DIR<br>
    </td>
  </tr>
  <tr>
    <td colspan="2">
      <b>Uninstall</b>
    </td>
  </tr>
  <tr>
    <td>
      brew uninstall git-town<br>
      brew untap Originate/gittown
    </td>
    <td>
      remove DIR<br>
      remove DIR from your path<br>
      remove DIR/man from your manpath
    </td>
  </tr>
</table>



## configuration

On first use, Git Town will ask for the main branch name and the names of any other non feature branches.
Git Town stores its configuration in the Git configuration of your project.
If these ever need to change, the configuration can be updated using <a href="http://git-scm.com/docs/git-config" target="_blank">git config</a>.


```bash
# Read configuration
git config git-town.main-branch-name
git config git-town.non-feature-branch-names

# Write configuration
git config git-town.main-branch-name master
git config git-town.non-feature-branch-names 'qa, production'
```

## documentation

In addition to the online documentation here,
you can run `git help town` on the command line
for an overview of the git town commands,
or `git help <command>` (e.g. `git help sync`)
for help on an individual command.


## development

tests are written in <a href="http://cukes.info/" target="_blank">Cucumber</a> and <a href="http://rspec.info/" target="_blank">RSpec</a>

```
# install tools
bundle
brew install shellcheck  # bash linter

# rake tasks
rake            # Run all linters and tests
rake lint       # Run all linters
rake lint:bash  # Run bash linter
rake lint:ruby  # Run ruby linter
rake spec       # Run tests

# run single test
cucumber -n 'scenario/feature name'
cucumber [filename]:[lineno]
```

Found a bug or want to contribute a feature?
<a href="https://github.com/Originate/git-town/issues/new" target="_blank">Open an issue</a>
or - even better - get down, go to town, and fire a feature-tested and linter-passing
<a href="https://help.github.com/articles/using-pull-requests" target="_blank">pull request</a>
our way!


## roadmap

The future roadmap is planned using
<a href="https://github.com/Originate/git-town/issues" target="_blank">GitHub issues</a>.
The past roadmap is in the <a href="release-notes.md" target="_blank">release notes</a>.

If you have an idea about a cool feature you would like to see in Git Town,
please <a href="https://github.com/Originate/git-town/issues/new" target="_blank">open a ticket</a>.
Our team will add the <a href="https://github.com/Originate/git-town/labels/idea" target="_blank">idea</a> tag.
Once we reach agreement about this idea, it will be tagged as <a href="https://github.com/Originate/git-town/labels/enhancement" target="_blank">enhancement</a>
or <a href="https://github.com/Originate/git-town/labels/bug" target="_blank">bug</a>.

