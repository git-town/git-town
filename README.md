# _Git Automation for Agile Development Teams_

<a href="https://travis-ci.org/Originate/git-town" alt="Build Status" target="_blank"><img src="https://travis-ci.org/Originate/git-town.svg?branch=master"></a>
[![License](http://img.shields.io/:license-MIT-blue.svg?style=flat)](MIT-LICENSE)

Git Town makes software development teams who use Git even more productive and happy.
It extends Git to support
[GitHub Flow](http://scottchacon.com/2011/08/31/github-flow.html),
[Git Flow](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow),
the [Nvie model](http://nvie.com/posts/a-successful-git-branching-model),
[GitLab Flow](https://about.gitlab.com/2014/09/29/gitlab-flow/),
and other workflows better, and allows to perform many common tasks faster and easier.

Check out [the big picture](documentation/big-picture.md) for more background on Git Town.


## Commands

Git Town provides these additional Git commands:

* [git extract](/documentation/git-extract.md) - copy selected commits from the current branch into their own branch
* [git hack](/documentation/git-hack.md) - cut a new feature branch off the main branch
* [git kill](/documentation/git-kill.md) - remove an obsolete feature branch
* [git pr](/documentation/git-pr.md) - create a new pull request
* [git prune-branches](/documentation/git-prune-branches.md) - delete merged branches
* [git repo](/documentation/git-repo.md) - view the repository homepage
* [git ship](/documentation/git-ship.md) - deliver a completed feature branch
* [git sync](/documentation/git-sync.md) - update the current branch with all relevant changes
* [git sync-fork](/documentation/git-sync-fork.md) - pull upstream updates into a forked repository
* [git town](/documentation/git-town.md) - general Git Town help, view and change Git Town configuration


## Installation

Git Town is written in 100% [Bash](https://www.gnu.org/software/bash/bash.html),
so it runs anywhere Git and Bash runs.

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
          (used by <a href="/documentation/git-extract.md">git extract</a>)
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
    mkdir -p ~/.config/fish/completions/
    curl -o ~/.config/fish/completions/git.fish http://raw.githubusercontent.com/Originate/git-town/master/autocomplete/git.fish
    ```


## Configuration

Git Town is configured on a per-repository basis.
Upon first use in a repository, Git Town will ask for all required
configuration.
Use the [git town](/documentation/git-town.md) command to view or update your configuration at any time.


## Documentation

In addition to the online documentation here,
you can run `git help town` on the command line
for an overview of the git town commands,
or `git help <command>` (e.g. `git help sync`)
for help on an individual command.


## Development

Tests are written in [Cucumber](http://cukes.info/) and [RSpec](http://rspec.info/).

```bash
# install tools needed for development
bundle
brew install shellcheck  # bash linter

# various ways to verify/test the code
rake            # Run all linters and specs
rake lint       # Run all linters
rake lint:bash  # Run bash linter
rake lint:ruby  # Run ruby linter
rake lint:cucumber  # Run cucumber linter
rake spec       # Run specs

# run a single test
cucumber [filename]:[lineno]

# run cucumber in parallel
bin/cuke [<folder>...]
```

Found a bug or want to contribute a feature?
Please [open an issue](https://github.com/Originate/git-town/issues/new)
or - even better - get down, go to town, and fire a feature-tested and linter-passing
[pull request](https://help.github.com/articles/using-pull-requests/)
our way!


## Roadmap

The future roadmap is planned using [GitHub issues](https://github.com/Originate/git-town/issues).
The past roadmap is in the [release notes](release-notes.md).

If you have an idea about a cool feature you would like to see in Git Town,
please [open a ticket](https://github.com/Originate/git-town/issues/new)!
