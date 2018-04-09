# Git Town Release Notes

## Unreleased

## 7.1.1 (2018-04-09)

#### Bug Fixes

* strip colors from the output of git commands run internally. This caused errors if you had git configured with `color.ui=always`

## 7.1.0 (2018-04-05)

#### New Features

* automatically remove outdated configuration

## 7.0.0 (2018-04-03)

#### BREAKING CHANGES

* `git town config`: `reset` and `setup` are now subcommands instead of flags
* `--abort`, `--continue`, `--skip`, `--undo` flags removed. Instead there are now top level commands `git town abort`, `git town continue`, `git town skip`, `git town undo`

#### New Features

* now catches when there is an unfinished state from a git town command that hit conflicts. If you try to run another git town command, you will be prompted on how to resolve the unfinished state. The unfinished state can be discarded and there is also a new top level command `git town discard` to delete the state of the last run command.

#### Bug Fixes

* skip perennial branch prompt if there are no options

## 6.0.2 (2018-01-26)

#### Bug Fixes

* fix parsing of git config when a value contains a newline

## 6.0.1 (2018-01-24)

#### Bug Fixes

* fixes displayed version number

## 6.0.0 (2018-01-15)

#### BREAKING CHANGES

* `git town set-parent-branch`: update interface to no longer accept arguments and instead prompt the user for the parent of the current branch
* `git town perennial-branches`: update the interface to add / remove perennial branches. Run `git town perennial-branch update` to receive the same prompt as initial configuration.
* Rename `hack-push-flag` to `new-branch-push-flag`. Please reconfigure if you are not using the default.

#### New Features

* `git town new-branch-push-flag`: add `--global` flag in order to set your default value. Any locally configured value will override.
* add `--debug` flag in order to see all the git commands runs under the hood
* speed improvement thanks to various optimizations to greatly reduce the number of git commands run under the hood

## 5.1.0 (2017-12-05)

#### New Features

* Nicer prompts from https://github.com/AlecAivazis/survey
* Parent branch prompt: add option to make the branch a perennial branch

#### Bug Fixes

* `git ship`: fix bug when encountering a merge conflict and using a code hosting driver ([#1060](https://github.com/Originate/git-town/issues/1060))

## 5.0.0 (2017-08-16)

#### BREAKING CHANGES

* `git new-pull-request / repo`: support for ssh identities changed
  * Previously: ssh identity needed to include "github", "gitlab" or "bitbucket"
  * Now: Run `git config git-town.code-hosting-origin-hostname <hostname>` where hostname matches what is in your ssh config file

#### New Features

* `git new-pull-request / repo`: support for self hosted versions
  * Run `git config git-town.code-hosting-driver <driver>` where driver is "bitbucket", "github", or "gitlab"
* `git sync`: add `--dry-run` flag to view the commands that would be run without running them
* `git ship`: when merging via the GitHub API, update the default commit message to include the PR title and number

## 4.2.1 (2017-08-16)

#### Bug Fixes

* add missing dependency to vendor folder (required for building on Homebrew)

## 4.2.0 (2017-08-15)

#### New Features

* Update all commands to support offline mode (lack of an internet connection)
  * Display / update offline mode with `git town offline [(true | false)]`
* `git ship`
  * add ability to ship hotfixes to perennial branches
  * add ability to merge via GitHub API when applicable. See [documentation](/documentation/commands/ship.md#github-pull-request-integration) for more info.

## 4.1.2 (2017-06-08)

#### Bug Fixes

* temporary file: use operating system temporary directory instead of hardcoding `/tmp`

## 4.1.1 (2017-06-07)

#### Bug Fixes

* temporary file: make parent directories if needed ([#955 comment](https://github.com/Originate/git-town/issues/955#issuecomment-306041043))

## 4.1.0 (2017-06-01)

#### New Features

* `git new-pull-request`, `git repo`: support more commands to open browsers (`cygstart`, `x-www-browser`, `firefox`, `opera`, `mozilla`, `netscape`)
* Add longer descriptions for commands which appear when running `git town help <command>` or `git town <command> --help`

#### Changes

* make `hack-push-flag` false by default (previously was true)
  ([#929](https://github.com/Originate/git-town/issues/929))

#### Bug Fixes

* replace all non-alpha numeric characters in temporary filename ([#925](https://github.com/Originate/git-town/issues/925))
* fix spacing in parent branch prompts
* enforce a minimum Git version of 2.7.0

## 4.0.1 (2017-05-21)

#### Bug Fixes

* fix infinite loop when prompting for parent branch and there are perennial branches configured
* enforce a minimum Git version of 2.6.0
* fix `ship` when the supplied branch is equal to the current branch and there are open changes
* allow `alias` to be run in a non-git directory

## 4.0.0 (2017-05-12)

#### BREAKING CHANGES

* rewrite in go, Git Town is now a single, stand-alone binary
  * first-class Windows support
  * This breaks existing aliases. If you have the default aliases setup,
    reconfigure them with `git town alias true`

## 3.1.0 (2017-03-27)

#### New Features

* `git new-pull-request`, `git repo`:
  * support `ssh://` urls (thanks to @zhangwei)
  * add GitLab support (thanks to @dgjnpr)

## 3.0.0 (2017-02-07)

#### BREAKING CHANGES

* `git hack`: no longer accepts a parent branch (functionality moved to `git append`)

#### New Features

* `git append`: create a new branch as a child of the current branch
* `git prepend`: create a new branch as a parent of the current branch
* `git rename-branch`: implicitly uses the current branch if only one branch name provided

#### Bug Fixes

* fix incorrectly reported branch loop
  ([#785](https://github.com/Originate/git-town/issues/785))

## 2.1.0 (2016-12-26)

#### New Features

* support multiple SSH identities
  ([#739](https://github.com/Originate/git-town/issues/739))

#### Bug Fixes

* update stashing strategy to avoid use of `git stash -u` which can delete ignored files
  ([#744](https://github.com/Originate/git-town/issues/744))
* fix merge conflicts resolution that results in no changes
  ([#753](https://github.com/Originate/git-town/issues/753))
* `git hack`: prompt for parent branch if unknown
  ([#760](https://github.com/Originate/git-town/issues/760))
* prevent parent branch loops
  ([#751](https://github.com/Originate/git-town/issues/751))

## 2.0.0 (2016-09-18)

#### BREAKING CHANGES

* All commands now have a `town-` prefix. Example `git town-sync`. This is to prevent conflicts with `git-extras` which adds git commands by the same name and `hub` which wants you to alias git to it and adds commands by the same name.
  * Use [git aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) to remove the `town-` prefix if you would like. Run `git town alias true` to add aliases for all `git-town` commands (skips commands which would overwrite existing aliases).

## 1.0.0 (2016-08-05)

#### New Features

* `git town set-parent-branch <child_branch_name> <parent_branch_name>`: to update a parent branch
  ([#729](https://github.com/Originate/git-town/issues/729))

#### Bug Fixes

* `git sync --all`: don't prompt for parent of perennial branches
  ([#727](https://github.com/Originate/git-town/issues/727))

## 0.10.1 (2016-06-23)

#### New Features

* `git hack`: add configuration for whether or not to push
  ([#720](https://github.com/Originate/git-town/issues/720))

#### Bug Fixes

* configuration: make branch order consistent
* `git ship`: update uncommitted changes error message
  ([#718](https://github.com/Originate/git-town/issues/718))

## 0.10.0 (2016-01-21)

#### BREAKING CHANGES

* `git prune-branches`: new functionality - delete branches whose tracking branch no longer exists
  ([#702](https://github.com/Originate/git-town/issues/702))

#### New Features

* `git sync`: add configuration option to merge the main branch / perennial branches with their upstream
  ([#671](https://github.com/Originate/git-town/issues/671))
* `git hack`, `git ship`: support for running in subfolders

#### Bug Fixes

* internationalize check for undefined function
  ([#678](https://github.com/Originate/git-town/issues/678))
* `git new-pull-request`: ability to continue after conflicts
  ([#700](https://github.com/Originate/git-town/issues/700))

## 0.9.0 (2015-10-17)

#### BREAKING CHANGES

* remove `git sync-fork`

#### New Features

* `git new-pull-request`: support forked repos
  ([#682](https://github.com/Originate/git-town/issues/682))
* `git sync`: if there is a remote upstream, the main branch is synced with its upstream counterpart
  ([#685](https://github.com/Originate/git-town/issues/685))

## 0.8.0 (2015-10-14)

#### BREAKING CHANGES

* removed `git extract`
* update internal storage of perennial branches
  * if you have configured multiple perennial branches, you will need to reset your configuration
    * `git town config --reset`
    * `git town config --setup` or you will be prompted the next time you run a Git Town command

#### New Features

* configuration prompt: allow user to select branch by number, ability to recover from bad input
* parent branch prompt: only show description and branch list once per command
* preserve checkout history so that `git checkout -` works as expected alongside Git Town commands
  ([#65](https://github.com/Originate/git-town/issues/65))
* `git hack`: pushes the new branch to the remote repository
  ([#664](https://github.com/Originate/git-town/issues/664))
* `git new-pull-request`: syncs the branch before creating the pull request
  ([#367](https://github.com/Originate/git-town/issues/367))
* `git sync --all`: pushes tags
  ([#464](https://github.com/Originate/git-town/issues/464))
* `git town config`: shows branch ancestry
  ([#651](https://github.com/Originate/git-town/issues/651))

#### Bug Fixes

* `git town version`: Homebrew installs no longer print date and SHA
  ([#631](https://github.com/Originate/git-town/issues/631))

## 0.7.3 (2015-09-02)

* `git kill`: remote only branch
  ([#380](https://github.com/Originate/git-town/issues/380))
* `git ship`: prompt when there are multiple authors
  ([#486](https://github.com/Originate/git-town/issues/486))

## 0.7.2 (2015-08-28)

* `git sync --all`: fix parent branch prompt
* `git ship`: comment out default commit message
  ([#382](https://github.com/Originate/git-town/issues/382))

## 0.7.1 (2015-08-27)

* `git ship`: add ability to ship remote only branch
  ([#593](https://github.com/Originate/git-town/issues/593))
* `git new-pull-request`/`git repo`: remove empty line output
  ([#602](https://github.com/Originate/git-town/issues/602))
* `git kill`: prompt for unknown parent branch
  ([#603](https://github.com/Originate/git-town/issues/603))
* `git sync --all`: prompt for unknown parent branch
  ([#604](https://github.com/Originate/git-town/issues/604))
* support branch names with forward slashes (along with any valid branch name)
  ([#608](https://github.com/Originate/git-town/issues/608))

## 0.7.0 (2015-08-24)

* fix `git ship --undo`
  ([#550](https://github.com/Originate/git-town/issues/550))
* rename `non-feature-branches` to `perennial-branches`
  ([#344](https://github.com/Originate/git-town/issues/344))
  * configuration is automatically updated to support this
* support for nested feature branches
  ([#529](https://github.com/Originate/git-town/issues/529))
* add `git rename-branch`
  ([#474](https://github.com/Originate/git-town/issues/474))
* rename `git pull-request` to `git new-pull-request`
  ([#413](https://github.com/Originate/git-town/issues/413), [#507](https://github.com/Originate/git-town/issues/507))
* add SHA and date to output of `git town version` for manual installs
* show error when trying to continue after a successful command
  ([#364](https://github.com/Originate/git-town/issues/364))

## 0.6.0 (2015-04-02)

* support for working without a remote repository for **git extract**, **git hack**, **git kill**, **git ship**, and **git sync**
  * implemented by our newest core committer @ricmatsui
* **git pr** renamed to **git pull-request**
  * set up an alias with `git config --global alias.pr pull-request`
* **git ship**
  * now accepts all `git commit` options
  * author with the most commits is automatically set as the author (when not the committer)
    ([#335](https://github.com/Originate/git-town/issues/335))
* **git pr/repo**
  * improved linux compatibility by trying `xdg-open` before `open`
* improved error messages when run outside a git repository
* improved setup wizard for initial configuration in a git repository
* added [contribution guide](/CONTRIBUTING.md)
* added tutorial

## 0.5.0 (2015-01-08)

* Manual installs need to update their `PATH` to point to the `src` folder within their clone of the repository
* **git extract:**
  * errors if branch exists remotely
    ([#236](https://github.com/Originate/git-town/issues/236))
  * removed restriction: need to be on a feature branch
    ([#269](https://github.com/Originate/git-town/issues/269))
  * added restriction: if no commits are provided, errors if the current branch does not have any have extractable commits (commits not in the main branch)
    ([#269](https://github.com/Originate/git-town/issues/269))
* **git hack:** errors if branch exists remotely
  ([#237](https://github.com/Originate/git-town/issues/237))
* **git kill:**
  * optional branch name
    ([#126](https://github.com/Originate/git-town/issues/126))
  * does not error if tracking branch has already been deleted
    ([#196](https://github.com/Originate/git-town/issues/196))
* **git pr:**
  * linux compatibility
    ([#232](https://github.com/Originate/git-town/issues/232))
  * compatible with more variants of specifying a Bitbucket or GitHub remote
    ([#271](https://github.com/Originate/git-town/issues/271))
  * compatible with respository names that contain ".git"
    ([#305](https://github.com/Originate/git-town/issues/305))
* **git repo:** view the repository homepage
  ([#140](https://github.com/Originate/git-town/issues/140))
* **git sync:**
  * `--all` option to sync all local branches
    ([#83](https://github.com/Originate/git-town/issues/83))
  * abort correctly after main branch updates and tracking branch conflicts
    ([#228](https://github.com/Originate/git-town/issues/228))
* **git town**: view and change Git Town configuration and easily view help page
  ([#98](https://github.com/Originate/git-town/issues/98))
* auto-completion for [Fish shell](http://fishshell.com)
  ([#177](https://github.com/Originate/git-town/issues/177))

## 0.4.1 (2014-12-02)

* **git pr:** create a new pull request
  ([#138](https://github.com/Originate/git-town/issues/138),
  [40d22e](https://github.com/Originate/git-town/commit/40d22eb1703ac96a58ac5052e70d20d7bdb9ac73))
* **git ship:**
  * empty commit message aborts the command
    ([#153](https://github.com/Originate/git-town/issues/153),
    [0bc84e](https://github.com/Originate/git-town/commit/0bc84ee626299896661fe1754cfa227630725bb9))
  * abort when there are no shippable changes
    ([#188](https://github.com/Originate/git-town/issues/188),
    [52fd94](https://github.com/Originate/git-town/commit/52fd94eca05bd3c2db5e7ac36121f08e56b9558b))
* **git sync:**
  * can now continue after just resolving conflicts (no need to commit or continue rebasing)
    ([#123](https://github.com/Originate/git-town/issues/123),
    [1a50ad](https://github.com/Originate/git-town/commit/1a50ad689a752f4eaed663e0ab22184621ee96a2))
  * restores deleted tracking branch
    ([#165](https://github.com/Originate/git-town/issues/165),
    [259464](https://github.com/Originate/git-town/commit/2594646ad853d83a6d697354d66755a374e42b8a))
* **git extract:** errors if branch already exists
  ([#128](https://github.com/Originate/git-town/issues/128),
  [75f498](https://github.com/Originate/git-town/commit/75f498771f19326f03bd1fd1bb70c9d9851b53f3))
* **git sync-fork:** no longer automatically sets upstream configuration
  ([865030](https://github.com/Originate/git-town/commit/8650301a3ea40a989562a991960fa0d41b26f7f7))
* remove needless checkouts for **git-ship**, **git-extract**, and **git-hack**
  ([#150](https://github.com/Originate/git-town/issues/150),
  [#155](https://github.com/Originate/git-town/issues/155),
  [8b385a](https://github.com/Originate/git-town/commit/8b385a745cf7ed28638e0a5c9c24440b7010354c),
  [35de43](https://github.com/Originate/git-town/commit/35de43156d9c6092840cd73456844b90acc36d8e))
* linters for shell scripts and ruby tests
  ([#149](https://github.com/Originate/git-town/issues/149),
  [076668](https://github.com/Originate/git-town/commit/07666825b5d60e15de274746fc3c26f72bd7aee2),
  [651c04](https://github.com/Originate/git-town/commit/651c0448309a376eee7d35659d8b06f709b113b5))
* rake tasks for development
  ([#170](https://github.com/Originate/git-town/issues/170),
  [ba74cf](https://github.com/Originate/git-town/commit/ba74cf30c8001941769dcd70410dbd18331f2fe9))

## 0.4.0 (2014-11-13)

* **git kill:** completely removes a feature branch
  ([#87](https://github.com/Originate/git-town/issues/87),
  [edd7d8](https://github.com/Originate/git-town/commit/edd7d8180eb76717fd72e77d2c75edf8e3b6b6ca))
* **git sync:** pushes tags to the remote when running on the main branch
  ([#68](https://github.com/Originate/git-town/issues/68),
  [71b607](https://github.com/Originate/git-town/commit/71b607988c00e6dfc8f2598e9b964cc2ed4cfc39))
* **non-feature branches:** cannot be shipped and do not merge main when syncing
  ([#45](https://github.com/Originate/git-town/issues/45),
  [31dce1](https://github.com/Originate/git-town/commit/31dce1dfaf11e1e17f17e141a26cb38360ab731a))
* **git ship:**
  * merges main into the feature branch before squash merging
    ([#61](https://github.com/Originate/git-town/issues/61),
    [82d4d3](https://github.com/Originate/git-town/commit/82d4d3e745732cb397850a4c047826ba485e2bdb))
  * errors if the feature branch is not ahead of main
    ([#86](https://github.com/Originate/git-town/issues/86),
    [a0ace5](https://github.com/Originate/git-town/commit/a0ace5bb5e992c193df8abe4b0aca984c302c323))
  * git ship takes an optional branch name
    ([#95](https://github.com/Originate/git-town/issues/95),
    [cbf020](https://github.com/Originate/git-town/commit/cbf020fc3dd6d0ce49f8814a92f103e243f9cd2b))
* updated output to show each git command and its output, updated error messages
  ([8d8973](https://github.com/Originate/git-town/commit/8d8973aaa58394a123ceed2811271606f4e1aaa9),
  [60e1d8](https://github.com/Originate/git-town/commit/60e1d8299ebbb0e75bdae057e864d17e1f9a3ce7),
  [408e69](https://github.com/Originate/git-town/commit/408e699e5bdd3af524b2ea64669b81fea3bbe60b))
* skips unnecessary pushes
  ([0da896](https://github.com/Originate/git-town/commit/0da8968aef29f9ecb7326e0fafb5976f51789dca))
* **man pages**
  ([609e11](https://github.com/Originate/git-town/commit/609e11400818604328885df86c02ee4630410e12),
  [164f06](https://github.com/Originate/git-town/commit/164f06bc8bf00d9e99ce0416f408cf62959dc833),
  [27b257](https://github.com/Originate/git-town/commit/27b2573ca5ffa9ae7930f8b5999bbfdd72bd16d9))
* **git prune-branches**
  ([#48](https://github.com/Originate/git-town/issues/48),
  [7a922e](https://github.com/Originate/git-town/commit/7a922ecd9e03d20ed5a0c159022e601cebc80313))
* **Cucumber:** optional Fuubar output
  ([7c5402](https://github.com/Originate/git-town/commit/7c540284cf46bd49a7623566c1343285813524c6))

## 0.3 (2014-10-10)

* multi-user support for feature branches
  ([#35](https://github.com/Originate/git-town/issues/35),
  [ca0882](https://github.com/Originate/git-town/commit/ca08820c68457bddf6b8fff6c3ef3d430b905d9b))
* **git sync-fork**
  ([#22](https://github.com/Originate/git-town/issues/22),
  [1f1f9f](https://github.com/Originate/git-town/commit/1f1f9f98ffa7288d6a5982ec0c9e571695590fe1))
* stores configuration in the Git configuration instead of a dedicated file
  ([8b8695](https://github.com/Originate/git-town/commit/8b86953d7c7c719f28dbc7af6e86d02adaf2053e))
* only makes one fetch from the central repo per session
  ([#15](https://github.com/Originate/git-town/issues/15),
  [43400a](https://github.com/Originate/git-town/commit/43400a5b968a47eb55484f73e34026f66b1e939a))
* automatically prunes remote branches when fetching updates
  ([86100f](https://github.com/Originate/git-town/commit/86100f08866f19a0f4e80f470fe8dcc6996ddc2c))
* always cleans up abort and continue scripts after using one of them
  ([3be4c0](https://github.com/Originate/git-town/commit/3be4c06635a943f378287963ba30e4306fcd9802))
* simpler readme, dedicated RDD document
* **<a href="http://cukes.info" target="_blank">Cucumber</a>** feature specs (you need Ruby 2.x)
  ([c9d175](https://github.com/Originate/git-town/commit/c9d175fe2f28fbda3f662454f54ed80306ce2f46))
* much faster testing thanks to completely local test Git repos
  ([#25](https://github.com/Originate/git-town/issues/25),
  [c9d175](https://github.com/Originate/git-town/commit/c9d175fe2f28fbda3f662454f54ed80306ce2f46))

## 0.2.2 (2014-06-10)

* fixes "unary" error messages
* lots of output and documentation improvements

## 0.2.1 (2014-05-31)

* better terminal output
* Travis CI improvements
* better documentation

## 0.2.0 (2014-05-29)

* displays the duration of specs
* only pulls the main branch if it has a remote
* --abort options to abort failed Git Town operations
* --continue options to continue some Git Town operations after fixing the underlying issues
* can be installed through Homebrew
* colored test output
* display summary after tests
* exit with proper status codes
* better documentation

## 0.1.0 (2014-05-22)

* git hack, git sync, git extract, git ship
* basic test framework
* Travis CI integration
* self-hosting: uses Git Town for Git Town development
