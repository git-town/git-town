![Git Town](https://raw.githubusercontent.com/git-town/git-town/master/website/img/git-town-horizontal.svg)

[![CircleCI](https://circleci.com/gh/git-town/git-town.svg?style=shield)](https://circleci.com/gh/git-town/git-town)
[![Go Report Card](https://goreportcard.com/badge/github.com/git-town/git-town)](https://goreportcard.com/report/github.com/git-town/git-town)
[![License](https://img.shields.io/:license-MIT-blue.svg?style=flat)](LICENSE)
[![Help Contribute to Open Source](https://www.codetriage.com/originate/git-town/badges/users.svg)](https://www.codetriage.com/originate/git-town)

Git Town makes [Git](https://git-scm.com) more efficient, especially for large
teams. See [this screencast](https://youtu.be/oLaUsUlFfTo) for an introduction
and this
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

**Git Town Installation**

- [git town alias](/documentation/commands/alias.md) - adds or removes shorter
  aliases for Git Town commands
- [git town install-fish-autocompletion](/documentation/commands/install-fish-autocompletion.md) -
  installs the autocompletion definition for [Fish shell](http://fishshell.com)
- [git town version](/documentation/commands/version.md) - displays the
  installed version of Git Town

## Installation

Since version 4.0, Git Town runs natively and without any dependencies on all
platforms. See the
[installation instructions](http://www.git-town.com/install.html) for more
details.

#### Aliasing

Git Town commands can be
[aliased](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) to remove the
`town` prefix:

```
git town alias true
```

After executing this, you can run `git hack` instead of `git town hack`. You can
also set this manually for individual commands:

```
git config --global alias.hack 'town hack'
```

## Configuration

Git Town prompts for required configuration information during usage. Use the
[git town config](/documentation/commands/config.md) command to manage the
stored configuration.

#### Required configuration

- the main development branch
- the
  [perennial branches](/documentation/development/branch_hierarchy.md#perennial-branches)

#### Optional Configuration

These configuration options have defaults, so the configuration wizard does not
ask about them.

- pull branch strategy

  - how to sync the main branch / perennial branches with their upstream
  - default: `rebase`
  - possible values: `merge`, `rebase`

- new branch push flag
  - whether or not branches created by hack / append / prepend should be pushed
    to remote repo
  - default: `false`
  - possible values: `true`, `false`

## Documentation

Run `git town` for an overview of the Git Town commands and `git help <command>`
(e.g. `git help sync`) for help with individual commands.

## Q&A

- **Is this compatible with my workflow?** <br> Yes. Git Town is compatible with
  [GitHub Flow](http://scottchacon.com/2011/08/31/github-flow.html),
  [Git Flow](https://www.atlassian.com/git/tutorials/comparing-workflows/feature-branch-workflow),
  the [Nvie model](https://nvie.com/posts/a-successful-git-branching-model),
  [GitLab Flow](https://about.gitlab.com/2014/09/29/gitlab-flow/), and most
  others workflows.

- **Does my whole team have to use Git Town?** <br> No. Just make sure that all
  feature branches get
  [squash-merged](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-request-merges#squash-and-merge-your-pull-request-commits),
  for example by requiring this in your
  [GitHub settings](https://help.github.com/en/github/administering-a-repository/configuring-commit-squashing-for-pull-requests).
  If you don't know what squash-merges are, you probably want to enable them.

## Contributing

Found a bug or have an idea for a new feature?
[Open an issue](https://github.com/git-town/git-town/issues/new) or send a
[pull request](https://help.github.com/articles/using-pull-requests/)! Check out
our [contributing](/CONTRIBUTING.md) and
[developer](documentation/development/README.md) guides to get started.

### Sponsors

Thanks to our sponsors for their continued support!

<table>
  <tr>
    <td>
      <a href="https://www.originate.com">
        <img src="documentation/originate.png" width="146" height="33">
      </a>
    </td>
  </tr>
</table>

### Contributors

Kudos to our contributors!

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
  </tr>
  <tr>
    <td align="center" width="60">
      <a href="https://github.com/iamandrewluca">
        <img src="https://avatars1.githubusercontent.com/u/1881266?s=60" width="60px">
        <sup><b>@iamandrewluca</b></sup>
      </a>
    </td>
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
      <a href="https://github.com/mjhm">
        <img src="https://avatars0.githubusercontent.com/u/431925?s=60" width="60px">
        <sup><b>@mjhm</b></sup>
      </a>
    </td>
  </tr>
  <tr>
    <td align="center" width="60">
      <a href="https://github.com/luketlancaster">
        <img src="https://avatars3.githubusercontent.com/u/8376505?s=60" width="60px">
        <sup><b>@luketlancaster</b></sup>
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
  <tr>
    <td align="center" width="60">
      <a href="https://github.com/WurmD">
        <img src="https://avatars2.githubusercontent.com/u/5755747?s=60" width="60px">
        <sup><b>@WurmD</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/sgarfinkel">
        <img src="https://avatars3.githubusercontent.com/u/10210461?s=60" width="60px">
        <sup><b>@sgarfinkel</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/cirego">
        <img src="https://avatars2.githubusercontent.com/u/551285?s=60" width="60px">
        <sup><b>@cirego</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/sheldonhull">
        <img src="https://avatars3.githubusercontent.com/u/3526320?s=60" width="60px">
        <sup><b>@sheldonhull</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/alphatroya">
        <img src="https://avatars0.githubusercontent.com/u/4927633?s=60" width="60px">
        <sup><b>@alphatroya</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/joshRpowell">
        <img src="https://avatars1.githubusercontent.com/u/6732638?s=60" width="60px">
        <sup><b>@joshRpowell</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/pattiereaves">
        <img src="https://avatars1.githubusercontent.com/u/44990?s=60" width="60px">
        <sup><b>@pattiereaves</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/zenspider">
        <img src="https://avatars0.githubusercontent.com/u/9832?s=60" width="60px">
        <sup><b>@zenspider</b></sup>
      </a>
    </td>
  </tr>
  <tr>
    <td align="center" width="60">
      <a href="https://github.com/Iron-Ham">
        <img src="https://avatars1.githubusercontent.com/u/3388381?s=60" width="60px">
        <sup><b>@Iron-Ham</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/grignaak">
        <img src="https://avatars2.githubusercontent.com/u/110779?s=60" width="60px">
        <sup><b>@grignaak</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/ericyliu">
        <img src="https://avatars2.githubusercontent.com/u/8580080?s=60" width="60px">
        <sup><b>@ericyliu</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/natecox">
        <img src="https://avatars0.githubusercontent.com/u/2782695?s=60" width="60px">
        <sup><b>@natecox</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/shadowhand">
        <img src="https://avatars3.githubusercontent.com/u/38203?s=60" width="60px">
        <sup><b>@shadowhand</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/hmbrg">
        <img src="https://avatars3.githubusercontent.com/u/7304269?s=60" width="60px">
        <sup><b>@hmbrg</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/qrevel">
        <img src="https://avatars2.githubusercontent.com/u/11804101?s=60" width="60px">
        <sup><b>@qrevel</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/aeneasr">
        <img src="https://avatars1.githubusercontent.com/u/3372410?s=60" width="60px">
        <sup><b>@aeneasr</b></sup>
      </a>
    </td>
  </tr>
  <tr>
    <td align="center" width="60">
      <a href="https://github.com/martinjaime">
        <img src="https://avatars1.githubusercontent.com/u/10568301?s=60" width="60px">
        <sup><b>@martinjaime</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/alexw10">
        <img src="https://avatars1.githubusercontent.com/u/9453636?s=60" width="60px">
        <sup><b>@alexw10</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/alexdavid">
        <img src="https://avatars2.githubusercontent.com/u/751581?s=60" width="60px">
        <sup><b>@alexdavid</b></sup>
      </a>
    </td>
    <td align="center" width="60">
      <a href="https://github.com/Braunson">
        <img src="https://avatars1.githubusercontent.com/u/577273?s=60" width="60px">
        <sup><b>@Braunson</b></sup>
      </a>
    </td>
  </tr>
</table>
