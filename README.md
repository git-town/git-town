![Git Town](https://git-town.github.io/git-town/documentation/logo-horizontal.svg)

[![Build Status](https://travis-ci.org/git-town/git-town.svg?branch=master)](https://travis-ci.org/git-town/git-town)
[![Go Report Card](https://goreportcard.com/badge/github.com/git-town/git-town)](https://goreportcard.com/report/github.com/git-town/git-town)
[![License](https://img.shields.io/:license-MIT-blue.svg?style=flat)](LICENSE)
[![Help Contribute to Open Source](https://www.codetriage.com/originate/git-town/badges/users.svg)](https://www.codetriage.com/originate/git-town)

Git Town makes software development teams who use Git even more productive and
happy. It adds Git commands that support
[GitHub Flow](http://scottchacon.com/2011/08/31/github-flow.html),
[Git Flow](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow),
the [Nvie model](https://nvie.com/posts/a-successful-git-branching-model),
[GitLab Flow](https://about.gitlab.com/2014/09/29/gitlab-flow/), and other
workflows more directly, and it allows you to perform many common Git operations
faster and easier.

See [git-town.com](https://www.git-town.com) for documentation and this
[Softpedia article](https://www.softpedia.com/get/Programming/Other-Programming-Files/Git-Town.shtml)
for an independent review.

## Commands

Git Town provides these additional Git commands:

**Development Workflow**

- [git town hack](/documentation/commands/hack.md) - cuts a new up-to-date
  feature branch off the main branch
- [git town sync](/documentation/commands/sync.md) - updates the current branch
  with all ongoing changes
- [git town new-pull-request](/documentation/commands/new-pull-request.md) -
  create a new pull request
- [git town ship](/documentation/commands/ship.md) - delivers a completed
  feature branch and removes it

**Repository Maintenance**

- [git town kill](/documentation/commands/kill.md) - removes a feature branch
- [git town prune-branches](/documentation/commands/prune-branches.md) - delete
  all merged branches
- [git town rename-branch](/documentation/commands/rename-branch.md) - rename a
  branch
- [git town append](/documentation/commands/append.md) - insert a new branch as
  a child of the current branch
- [git town prepend](/documentation/commands/prepend.md) - insert a new branch
  between the current branch and its parent
- [git town repo](/documentation/commands/repo.md) - view the repository
  homepage

**Git Town Configuration**

- [git town config](/documentation/commands/config.md) - displays or updates
  your Git Town configuration
- [git town new-branch-push-flag](/documentation/commands/new-branch-push-flag.md) -
  configures whether new empty branches are pushed to origin
- [git town main-branch](/documentation/commands/main-branch.md) - displays or
  sets the main development branch for the current repo
- [git town offline](/documentation/commands/offline.md) - enables/disables
  offline mode
- [git town perennial-branches](/documentation/commands/perennial-branches.md) -
  displays or updates the perennial branches for the current repo
- [git town pull-branch-strategy](/documentation/commands/pull-branch-strategy.md) -
  displays or sets the strategy with which perennial branches are updated
- [git town set-parent-branch](/documentation/commands/set-parent-branch.md) -
  updates a branch's parent

**Other Commands**

- [git town alias](/documentation/commands/alias.md) - adds or removes shorter
  aliases for Git Town commands
- [git town install-fish-autocompletion](/documentation/commands/install-fish-autocompletion.md) -
  installs the autocompletion definition for [Fish shell](http://fishshell.com)
- [git town version](/documentation/commands/version.md) - displays the
  installed version of Git Town

## Installation

Since version 4.0, Git Town runs natively on all platforms without any
dependencies. Check out our
[installation instructions](http://www.git-town.com/install.html) for more
details.

### Aliasing

Each command can be
[aliased](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) individually to
remove the `town` prefix with:

```
git config --global alias.hack 'town hack'
```

Now you can run `git hack` instead of `git town hack`. As a convenience, you can
add or remove global aliases for all `git-town` commands with:

```
git town alias <true | false>
```

## Configuration

Git Town is configured on a per-repository basis. Upon first use in a
repository, you will be prompted for the required configuration. Use the
[git town config](/documentation/commands/config.md) command to view or update
your configuration at any time.

#### Required configuration

- the main development branch
- the
  [perennial branches](/documentation/development/branch_hierarchy.md#perennial-branches)

#### Optional Configuration

The following configuration options have defaults, so the configuration wizard
does not ask about them.

- the pull branch strategy

  - how to sync the main branch / perennial branches with their upstream
  - default: `rebase`
  - possible values: `merge`, `rebase`

- the new branch push flag
  - whether or not branches created by hack / append / prepend should be pushed
    to remote repo
  - default: `false`
  - possible values: `true`, `false`

## Documentation

In addition to the online documentation here, you can run `git town` on the
command line for an overview of the Git Town commands, or `git help <command>`
(e.g. `git help sync`) for help with an individual command.

## Contributing

Thanks to our sponsors for their continued support!

<table>
  <tr>
    <td>
      <a href="https://www.originate.com" target="_blank">
        <img src="documentation/originate.png" width="146" height="33">
      </a>
    </td>
  </tr>
</table>

A big thank to our contributors!

<table>
  <tr>
    <td align="center" width="60">
      <a href="https://github.com/kevgo">
        <img src="https://avatars.githubusercontent.com/u/268934?s=60" width="60px">
        <sup><b>@kevgo</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/charlierudolph">
        <img src="https://avatars2.githubusercontent.com/u/1676758?s=60" width="60px">
        <sup><b>@charlierudolph</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/allewun">
        <img src="https://avatars2.githubusercontent.com/u/1256911?s=60" width="60px">
        <sup><b>@allewun</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/ricmatsui">
        <img src="https://avatars0.githubusercontent.com/u/5288285?s=60" width="60px">
        <sup><b>@ricmatsui</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/redouglas">
        <img src="https://avatars3.githubusercontent.com/u/1149609?s=60" width="60px">
        <sup><b>@redouglas</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/allonsy">
        <img src="https://avatars2.githubusercontent.com/u/5892756?s=60" width="60px">
        <sup><b>@allonsy</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/aricahunter">
        <img src="https://avatars2.githubusercontent.com/u/5395515?s=60" width="60px">
        <sup><b>@aricahunter</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/jiexi">
        <img src="https://avatars2.githubusercontent.com/u/918701?s=60" width="60px">
        <sup><b>@jiexi</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/iamandrewluca">
        <img src="https://avatars1.githubusercontent.com/u/1881266?s=60" width="60px">
        <sup><b>@iamandrewluca</b></sup>
      </a>
    </td>
  </tr>
  <tr>
    <td align="center" width="60">
      <a href="https://github.com/ChrisMM">
        <img src="https://avatars0.githubusercontent.com/u/1051777?s=60" width="60px">
        <sup><b>@ChrisMM</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/zhangwei">
        <img src="https://avatars2.githubusercontent.com/u/6028709?s=60" width="60px">
        <sup><b>@zhangwei</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/tallpants">
        <img src="https://avatars2.githubusercontent.com/u/15325890?s=60" width="60px">
        <sup><b>@tallpants</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/seanstrom">
        <img src="https://avatars3.githubusercontent.com/u/2845768?s=60" width="60px">
        <sup><b>@seanstrom</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/schneems">
        <img src="https://avatars2.githubusercontent.com/u/59744?s=60" width="60px">
        <sup><b>@schneems</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/morrme">
        <img src="https://avatars1.githubusercontent.com/u/26514778?s=60" width="60px">
        <sup><b>@morrme</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/">
        <img src="" width="60px">
        <sup><b>@</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/mjhm">
        <img src="https://avatars0.githubusercontent.com/u/431925?s=60" width="60px">
        <sup><b>@mjhm</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/luketlancaster">
        <img src="https://avatars3.githubusercontent.com/u/8376505?s=60" width="60px">
        <sup><b>@luketlancaster</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/">
        <img src="" width="60px">
        <sup><b>@</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/ecbrodie">
        <img src="https://avatars1.githubusercontent.com/u/1844664?s=60" width="60px">
        <sup><b>@ecbrodie</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/doismellburning">
        <img src="https://avatars1.githubusercontent.com/u/817118?s=60" width="60px">
        <sup><b>@doismellburning</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/dgjnpr">
        <img src="https://avatars1.githubusercontent.com/u/1767441?s=60" width="60px">
        <sup><b>@dgjnpr</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/">
        <img src="" width="60px">
        <sup><b>@</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/atilacamurca">
        <img src="https://avatars1.githubusercontent.com/u/508624?s=60" width="60px">
        <sup><b>@atilacamurca</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/Zearin">
        <img src="https://avatars0.githubusercontent.com/u/630124?s=60" width="60px">
        <sup><b>@Zearin</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/TKAB">
        <img src="https://avatars3.githubusercontent.com/u/66597?s=60" width="60px">
        <sup><b>@TKAB</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/Siilwyn">
        <img src="https://avatars2.githubusercontent.com/u/5701149?s=60" width="60px">
        <sup><b>@Siilwyn</b></sup>
      </a>
    </td>
  </tr>
</table>

Found a bug or have an idea for a new feature?
[Open an issue](https://github.com/git-town/git-town/issues/new) or - even
better - get down, go to town, and fire a feature-tested
[pull request](https://help.github.com/articles/using-pull-requests/) our way!
Check out our [contributing guide](/CONTRIBUTING.md) to start coding.
