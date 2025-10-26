# Git Town Changelog

## 22.2.0 (unreleased)

#### New Features

- the new display-types configuration setting allows configuring which branch types Git Town shows in lists of branches. This affects the [branch](https://www.git-town.com/commands/branch.html), [switch](https://www.git-town.com/commands/switch.html), [config](https://www.git-town.com/commands/config.html), [down](https://www.git-town.com/commands/down.html), [up](https://www.git-town.com/commands/up.html) commands and all internal dialogs that ask to select a branch or parent branch.

## 22.1.0 (2025-10-13)

#### New Features

- if you run a Git Town command while another is suspended, Git Town now offers the option to finish the suspended program and then run the new one ([#3337](https://github.com/git-town/git-town/issues/3337)).

#### Bug Fixes

- Fixes bugs resulting from the new Option serialization ([#5623](https://github.com/git-town/git-town/discussions/5623)).

#### Contributors

Shoutout to @DPirate, @IGassmann, @Mause, @derekspelledcorrectly, @ethankeshishian, @kevgo, @stephenwade, @yaadata, @Shmookoff for contributing code, bug reports, and ideas to 24 shipped pull requests and 4 resolved issues!

## 22.0.0 (2025-09-23)

#### BREAKING CHANGES

- Renamed the `codeberg` connector to `forgejo` since it supports all Forgejo-based forges. Codeberg itself runs on Forgejo ([#5447](https://github.com/git-town/git-town/issues/5447)).
- Start the setup assistant with `git town init` instead of `git town config setup`, matching Git's `git init` ([#5269](https://github.com/git-town/git-town/issues/5269)).
- [git town hack](https://www.git-town.com/commands/hack.html) no longer converts an existing branch into a feature branch. Use the new [feature](https://www.git-town.com/commands/feature.html) command instead ([#5516](https://github.com/git-town/git-town/pull/5516)).
- The `contribute`, `feature`, `observe`, `park`, and `prototype` commands no longer signal a problem when you run them on a branch that already has the desired type.
- Removed the long-deprecated `kill` command. Use [delete](https://www.git-town.com/commands/delete.html) instead.

#### New Features

- Added support for Azure DevOps ([#1657](https://github.com/git-town/git-town/issues/1657)).
- Introduced the [push-branches](https://www.git-town.com/preferences/push-branches.html) setting to control whether Git Town pushes local changes to tracking branches automatically. Disable it if you prefer to push manually ([#5541](https://github.com/git-town/git-town/issues/5541)).
- Introduced the [auto-sync](https://www.git-town.com/preferences/auto-sync.html) setting to control whether [hack](https://www.git-town.com/commands/hack.html), [append](https://www.git-town.com/commands/append.html), and [prepend](https://www.git-town.com/commands/prepend.html) sync existing branches before creating a new one ([#5540](https://github.com/git-town/git-town/issues/5540)).
- `set-parent`: added the `--none` flag to create perennial branches programmatically.

#### Bug Fixes

- Fixed `git town branch` in repos with a detached HEAD ([#5565](https://github.com/git-town/git-town/issues/5565)).

#### Contributors

Shoutout to @Phunky, @alexus37, @aximut, @bb010g, @benmosher, @blaggacao, @charlierudolph, @derekspelledcorrectly, @kevgo, @kinyat, @maruffahmed, @oludaara, @pradeepmurugesan, @sheldonhull, @stefanfrede, @stephenwade, @vectro, @whitebear-gh, @yaadata for contributiong code, feedback, and ideas to 78 shipped pull requests and 12 resolved issues!

## 21.5.0 (2025-09-05)

#### New Features

- Added [up](https://www.git-town.com/commands/up.html) and [down](https://www.git-town.com/commands/down.html) commands to quickly move to the child or parent of the current branch ([#5432](https://github.com/git-town/git-town/issues/5432)).
- All configuration settings can now be provided via environment variables, in addition to the config file and Git metadata. This makes it easy to use custom scripts that provide configuration data. For example, to load an API token from 1password CLI:

  ```bash
  GIT_TOWN_GITHUB_TOKEN=$(op read op://development/GitHub/credentials/personal_token) git town config
  ```

  ([#5446](https://github.com/git-town/git-town/issues/5446))
- The `hack`, `append`, and `prepend` commands now support a [stash](https://www.git-town.com/commands/hack.html#--stash--no-stash) flag and config option to leave staged changes as-is. Handy if you've carefully prepared an index you want to commit to a new branch ([#5429](https://github.com/git-town/git-town/issues/5429)).
- The setup assistant can now do a quick setup with only the essential settings ([#5484](https://github.com/git-town/git-town/issues/5484)).
- Added a new `feature` command to convert the given or current branch into a feature branch ([#5376](https://github.com/git-town/git-town/issues/5376)).
- The [detached sync](https://www.git-town.com/commands/sync.html#-d--detached) flag can now be permanently enabled through the new [detached configuration setting](https://www.git-town.com/preferences/detached.html) ([#5452](https://github.com/git-town/git-town/issues/5452)).
- The [auto-resolve](https://www.git-town.com/commands/sync.html#--auto-resolve) flag can now be disabled with `--no-auto-resolve` ([#5458](https://github.com/git-town/git-town/pull/5458)).

#### Bug Fixes

- The setup assistant no longer stores a dev-remote when the user selects the default option ([#5492](https://github.com/git-town/git-town/issues/5492)).

#### Contributors

Shoutout to @fuadsaud, @kevgo, @stefanfrede, @stephenwade, @yaadata for contributing code, feedback, and ideas to 65 shipped PRs and 10 resolved issues!

## 21.4.3 (2025-08-15)

#### Bug Fixes

- Fewer phantom merge conflicts: Git Town now performs a rebase-onto only if there are actual commits to remove. If there are no commits to remove, Git Town performs a regular rebase, or if there is no need to sync, no rebase at all. ([#5422](https://github.com/git-town/git-town/pull/5422))
- "git town branch" prints branches in other worktrees de-emphasized ([#5405](https://github.com/git-town/git-town/pull/5405))

#### Contributors

Shoutout to @AmitJoki, @Ydot19, @avaz, @benmosher, @kevgo, @nebbles, @nekitk, @stephenwade, @tranhl for contributing feedback, ideas, and code to 47 shipped PRs and 3 resolved issues!

## 21.4.2 (2025-08-11)

#### Bug Fixes

- Git Town no longer mistakes legit file conflicts for phantom conflicts ([#5156](https://github.com/git-town/git-town/issues/5156), [#5140](https://github.com/git-town/git-town/issues/5140))
- delete: Git Town now rebases onto the correct branch ([#5358](https://github.com/git-town/git-town/pull/5358))

#### Contributors

Shoutout to @AmitJoki, @Ydot19, @avaz, @benmosher, @ceilfors, @kevgo, @nebbles, @nekitk, @stephenwade, @tim-richter, @towry, @tranhl for contributing detailed bug reproductions and ideas to 55 shipped PRs and 4 resolved issues!

## 21.4.1 (2025-08-05)

#### Bug Fixes

- Fixes parsing of the new Git metadata config entries that disable auto-resolving phantom conflicts ([#5326](https://github.com/git-town/git-town/pull/5326)).

## 21.4.0 (2025-08-01)

#### New Features

- The setup assistant now correctly handles existing configuration settings in global Git metadata ([#5201](https://github.com/git-town/git-town/issues/5201))
- When Git Town is unconfigured, it now runs the full setup assistant instead of just asking for the main branch ([#5057](https://github.com/git-town/git-town/issues/5057))
- `git town set-parent` now presents the same hierarchical branch list as `git town switch` ([#5259](https://github.com/git-town/git-town/pull/5259))
- When prompting for a parent branch, Git Town now also shows the hierarchical list from `git town switch` ([#5266](https://github.com/git-town/git-town/pull/5266))
- Added a configuration option and CLI switch to disable automatic resolving of phantom merge conflicts for cases where manual conflict resolution is preferred ([#5317](https://github.com/git-town/git-town/pull/5317))

#### Bug Fixes

- Phantom merge conflicts are now auto-resolved even when your stack is rooted in a perennial branch instead of the main branch ([#5193](https://github.com/git-town/git-town/issues/5193))
- The `glab` connector now correctly updates merge proposal targets ([#5283](https://github.com/git-town/git-town/pull/5283))
- Git Town now also removes branch type overrides when it deletes branches that were shipped at the remote ([#5274](https://github.com/git-town/git-town/pull/5274))
- The setup assistant now cleans up local Git metadata when you opt to save configuration globally ([#5230](https://github.com/git-town/git-town/pull/5230))

#### Contributors

Shoutout to @Ydot19, @kevgo, @thekarel for contributing code, ideas, and bug reports to 145 shipped PRs 5 resolved issues! 🚀

## 21.3.0 (2025-07-16)

#### New Features

_[setup assistant](https://www.git-town.com/commands/config-setup.html):_

- now creates configuration files with name `git-town.toml` instead of the deprecated `git-branches.toml` ([#5162](https://github.com/git-town/git-town/pull/5162))
- now gives the user a chance to enter the [observed regex](https://www.git-town.com/preferences/observed-regex.html) and [contribution regex](https://www.git-town.com/preferences/contribution-regex.html), ([#5133](https://github.com/git-town/git-town/pull/5133), [#5132](https://github.com/git-town/git-town/pull/5132))
- when entering perennial branches, the main and perennial branches from the config file are now preselected and locked ([#5154](https://github.com/git-town/git-town/pull/5154))
- no longer asks for the [dev remote](https://www.git-town.com/preferences/dev-remote.html) if only one Git remote exists ([#5153](https://github.com/git-town/git-town/pull/5153))
- now updates Git metadata only if the user entered a different value than already exists ([#5127](https://github.com/git-town/git-town/pull/5127))
- no longer creates commented out entries ([#5110](https://github.com/git-town/git-town/pull/5110))
- now explains how to run it manually ([#5155](https://github.com/git-town/git-town/pull/5155/files))

#### Bug Fixes

- updating the base branch in a stack using the [gh connector](https://www.git-town.com/preferences/github-connector.html) works now ([#5163](https://github.com/git-town/git-town/pull/5163))
- setup assistant:
  - now displays the entered value for the API token scope ([#5144](https://github.com/git-town/git-town/pull/5144))
  - more consistent dialog captions ([#5159](https://github.com/git-town/git-town/pull/5159))

#### Contributors

Shoutout to @kevgo, @stephenwade, @wengh for contributing code, ideas, and feedback to 52 shipped PRs and 1 resolved issue!

## 21.2.0 (2025-07-02)

#### New Features

- Git Town can now use GitHub's [gh CLI](https://cli.github.com) to talk to the GitHub API. No more messing with access tokens manually! ([#1639](https://github.com/git-town/git-town/issues/1639))
- GitLab users get the same treatment: Git Town now integrates with the [glab CLI](https://gitlab.com/gitlab-org/cli/-/tree/main) to access the GitLab API ([#5079](https://github.com/git-town/git-town/pull/5079)).
- Dropped the dependency on `which` on Unix-like systems. One less external tool to worry about ([#5060](https://github.com/git-town/git-town/pull/5060)).
- Git Town is now in the official Arch Linux repositories. Install with `pacman -S git-town` ([#5015](https://github.com/git-town/git-town/pull/5015)).
- Git Town is now also available on OpenSUSE [Stable](https://github.com/git-town/git-town/pull/5058) and [Tumbleweed](https://github.com/git-town/git-town/issues/5032).
- The Setup Assistant now validates the forge information you enter works. If the connection fails, you get a chance to enter the credentials again. No more silent misconfigurations. ([#3030](https://github.com/git-town/git-town/issues/3030)).

#### Bug Fixes

- `git town diff-parent` now shows only the changes introduced by the current branch, even when it's behind its parent ([#5053](https://github.com/git-town/git-town/pull/5053)).
- The GitLab connector now handles the `--title` and `--body` arguments correctly ([#5072](https://github.com/git-town/git-town/issues/5072)).
- Continuing a suspended Git Town command that needs forge access now works correctly in all edge cases ([#5098](https://github.com/git-town/git-town/pull/5098)).
- All interactive dialogs now render properly in 80-character-wide terminals ([#5074](https://github.com/git-town/git-town/pull/5074)).
- The setup assistant now correctly pre-selects the token scope you previously configured ([#5046](https://github.com/git-town/git-town/pull/5046)).

#### Contributors

Huge thanks to @ChrisMM, @JafethAriasH, @alerque, @alphatroya, @ccoVeille, @emmanuel-ferdman, @haltcase, @kastl-ars, @kelbyers, @kevgo, @stephenwade, @tranhl, @vectro, @znd4 for contributing ideas, feedback, code, and installer support to 86 shipped PRs and 9 resolved issues. Cheers!

## 21.1.0 (2025-06-04)

#### New Features

- Git Town now keeps an immutable, append-only log of the repository state before and after each Git Town operation. This is a safety net in case `git town undo` cannot undo all changes. The new `git town runlog` command prints the log as well as its file path ([#4456](https://github.com/git-town/git-town/issues/4456)).

#### Bug Fixes

- Git Town now supports the situation where a file has the same name as a branch ([#5001](https://github.com/git-town/git-town/issues/5001)).
- Fixes a bug where `git town sync` sometimes doesn't push local changes ([#5007](https://github.com/git-town/git-town/issues/5007)).

#### Contributors

Shoutout to @AmitJoki, @SmolPandaDev, @kevgo, @legeana, @niklashigi, @stephenwade, @tobiaseisenschenk, @towry, @tranhl for contributing valuable code, feedback, and ideas to 19 shipped PRs and 6 resolved issues!

## 21.0.0 (2025-05-27)

#### BREAKING CHANGES

- **Configuration setting `default-branch-type` is now [unknown-branch-type](https://www.git-town.com/preferences/unknown-branch-type.html).** This better reflects that this setting applies to branches without a known type, and helps differentiate it from [new-branch-type](https://www.git-town.com/preferences/new-branch-type.html). Existing configs continue to work indefinitely. Git-based configuration gets updated automatically, updating this in the config file is recommended ([#4964](https://github.com/git-town/git-town/issues/4964)).
- **Updated branch name during merge.** When [merging](https://www.git-town.com/commands/merge.html) two branches, Git Town now uses the parent branch for the name of the merged branch instead of the child branch. This keeps the pull request of the parent branch intact and generally aligns better with typical usage of this command ([#4938](https://github.com/git-town/git-town/issues/4938)).
- **[create.new-branch-type](https://www.git-town.com/preferences/new-branch-type.html) is now always respected.** Previously, if this config option was set to `feature`, Git Town didn't apply it. Now it always assigns the configured branch type ([#4946](https://github.com/git-town/git-town/issues/4946)).

#### New Features

- **New [walk](https://www.git-town.com/commands/walk.html) command:** Execute a shell command on all branches in a stack or your workspace. Without a command it exits to the shell for each branch. Great for applying automated changes to all branches or debugging issues like which branch breaks a linter ([#4852](https://github.com/git-town/git-town/issues/4852)).
- **Smarter syncing:** `git town sync` now skips Git operations that wouldn't result in any changes. This speeds things up and avoids unnecessary Git noise ([#4927](https://github.com/git-town/git-town/issues/4927)).

#### Bug Fixes

- Beamed commits are now always removed from their original location after being moved ([#4895](https://github.com/git-town/git-town/issues/4895)).
- More reliable detection of the first commit in a branch, reducing edge case failures ([#4980](https://github.com/git-town/git-town/issues/4980))
- `git town branch` no longer shows duplicate branches when multiple Git remotes are present ([#4961](https://github.com/git-town/git-town/issues/4961)).

#### Contributors

Huge thanks to @AmitJoki, @WhosNickDoglio, @jfmyers9, @kevgo, @mw00120, @ruudk, @stephenwade, @zodman for moving Git Town forward by contributing code, feedback, and ideas to 52 shipped PRs and 13 resolved issues!

## 20.2.0 (2025-05-15)

#### New Features

- `git town sync` now only executes the Git operations that are actually needed, i.e. skips Git operations that would do nothing ([#4913](https://github.com/git-town/git-town/pull/4913), [#4907](https://github.com/git-town/git-town/pull/4907), [#4902](https://github.com/git-town/git-town/pull/4902)).
- When shipping via the forge API, `git town ship` now pre-populates the commit message with the proposal title and description ([#2095](https://github.com/git-town/git-town/issues/2095)).
- Git Town now supports detecting copied files via the `git config diff.renames copies` setting ([#4878](https://github.com/git-town/git-town/issues/4878)).

#### Bug Fixes

- Detached mode no longer pulls updates from the local main branch into feature branches ([#4890](https://github.com/git-town/git-town/issues/4890)).
- Git Town no longer deletes details when updating pull requests on BitBucket Cloud ([#4900](https://github.com/git-town/git-town/issues/4900)).
- `git town sync` with the `compress` strategy no longer re-creates commits if there are no changes ([#4342](https://github.com/git-town/git-town/issues/4342)).
- `git town sync` now respects the `--no-push` setting when rebasing ([#4930](https://github.com/git-town/git-town/issues/4930)).
- Git Town no longer unstashes if the initial stash command was ineffectual ([#1003](https://github.com/git-town/git-town/issues/1003)).
- `git town propose` now always runs detached ([#4915](https://github.com/git-town/git-town/pull/4915)).
- `git town sync` now more reliably skips the editor ([#4911](https://github.com/git-town/git-town/pull/4911)).

#### Contributors

Shoutout to @AmitJoki, @JCB-K, @blaggacao, @charlierudolph, @erik-rw, @fcurella, @kevgo, @legeana, @mw00120, @sheldonhull, @stephenwade for contributing code, feedback, and ideas to 59 shipped PRs and 13 resolved issues!

## 20.1.0 (2025-05-08)

#### New Features

- `git town compress` now has a `--no-verify` flag that disables Git's pre-commit hook ([#4843](https://github.com/git-town/git-town/issues/4843)).

#### Bug Fixes

- `git town compress` now enforces that the branch to compress is in sync with its parent branch ([#4845](https://github.com/git-town/git-town/issues/4845)).
- `git town sync` now doesn't remove commits of branches with deleted tracking branch if they don't have descendents ([#4872](https://github.com/git-town/git-town/discussions/4872)).
- Git Town no longer overrides the language of executed Git commands to US-English ([#4861](https://github.com/git-town/git-town/pull/4861)).

#### Contributors

Shoutout to @AmitJoki, @fcurella, @haltcase, @kevgo, @lvlcn-t, @mw00120, @niklastreml, @stephenwade for contributing code, feedback, and ideas to 34 shipped PRs and 6 resolved issues!

## 20.0.0 (2025-05-02)

Git Town 2000! 🎉

#### BREAKING CHANGES

- The `push-new-branches` configuration option is now called `share-new-branches` and allows additional ways of sharing freshly created branches ([#3912](https://github.com/git-town/git-town/issues/3912)):
  - `no`: keep new branches local (default)
  - `push`: push new branches to the [development remote](https://www.git-town.com/preferences/dev-remote.html)
  - `propose`: automatically create proposals for new branches. This helps being maximally transparent with progress on each item worked on.
- `git town propose` now always syncs the proposed branch, but always in [detached mode](https://www.git-town.com/commands/sync.html#-d--detached) mode ([#4772](https://github.com/git-town/git-town/pull/4772), [#4781](https://github.com/git-town/git-town/issues/4781)).
- `git town propose` now longer has the `--detached` flag because it now always syncs in detached mode ([#4775](https://github.com/git-town/git-town/pull/4775)).

#### New Features

- `git town sync` now correcly syncs branches whose commits got amended or rebased ([#4586](https://github.com/git-town/git-town/issues/4586)).
- You can now propose all branches in a stack with `git town propose --stack` ([#3840](https://github.com/git-town/git-town/issues/3840)).
- `git town propose` now un-parks parked branches when proposing them ([#4780](https://github.com/git-town/git-town/issues/4780)).
- The [setup assistant](https://www.git-town.com/configuration.html) no longer asks for data already provided by the [Git Town configuration file](https://www.git-town.com/configuration-file.html) ([#4710](https://github.com/git-town/git-town/issues/4710)).
- The setup assistant now offers to store forge API tokens globally for all repos on your machine ([#4112](https://github.com/git-town/git-town/issues/4112)).
- [git town status reset](https://www.git-town.com/commands/status-reset.html) now indicates whether the runstate file existed ([#4814](https://github.com/git-town/git-town/pull/4814)).
- [git town status reset](https://www.git-town.com/commands/status-reset.html) now supports the `--verbose` flag ([#4813](https://github.com/git-town/git-town/pull/4813)).

#### Bug Fixes

- Git Town now correctly resolves `includeIf` directives in Git configuration ([#4107](https://github.com/git-town/git-town/issues/4107)).
- `git town prepend --beam` now works correctly with prototype branches ([#4768](https://github.com/git-town/git-town/issues/4768)).
- Git Town now loads the forge API token with the same precendence as other configuration data ([#7428](https://github.com/git-town/git-town/pull/4728)).
- [git town undo](https://www.git-town.com/commands/undo.html) now correctly undoes situations where only the local part of a branch got renamed ([#4794](https://github.com/git-town/git-town/pull/4794)).
- Git Town now works even if Git's `color.ui` setting is `always` ([#4840](https://github.com/git-town/git-town/pull/4840)).
- The setup assistant now only updates the stored access token of the forge that is actually being used ([#4819](https://github.com/git-town/git-town/issues/4819)).
- [git town status reset](https://www.git-town.com/commands/status-reset.html) can now be run from a subdirectory ([#4812](https://github.com/git-town/git-town/pull/4812)).

#### Contributors

Git Town 2000 is a big release. Shoutout to @AmitJoki, @Ydot19, @ahgraber, @davidolrik, @erik-rw, @haltcase, @jmyers-figma, @judec-ps, @kevgo, @lvlcn-t, @nekitk, @niklastreml, @pradeepmurugesan, @ruudk, @stephenwade, @terheyden, @tharun208 for contributing code, feedback, and ideas to 124 shipped PRs and 17 resolved issues!

## 19.0.0 (2025-04-17)

#### BREAKING CHANGES

- The commands `new-pull-request` and `rename-branch` are being sunset after being deprecated for a long time. Their modern replacements are `propose` and `rename` ([#4714](https://github.com/git-town/git-town/pull/4714)).
- The configuration entries `contribution-branches`, `observed-branches`, `parked-branches`, and `prototype-branches` are being sunset. Their functionality is taken over by setting the type for individual branches as well as `contribution-regex`, `observed-regex`, `default-branch-type`, and `new-branch-type` ([#4499](https://github.com/git-town/git-town/issues/4499)).

#### New Features

- `git town append` and `git town hack` now also have a `--beam` flag to move selected commits to the new branch. When enabled, they no longer fetch or sync, which allows you to move commits with the fewest possible distractions ([#3338](https://github.com/git-town/git-town/issues/3338)).
- The "select commits to beam" dialog now displays the SHA of commits in addition to the commit message ([#4519](https://github.com/git-town/git-town/issues/4519)).
- `set-parent` now allows providing the new parent as an optional positional CLI argument ([documentation](https://www.git-town.com/commands/set-parent.html#positional-argument), [#4705](https://github.com/git-town/git-town/issues/4705)).
- The Git Town website now has a [how-to](https://www.git-town.com/how-tos.html) section.

#### Bug Fixes

- `git town sync --no-push` no longer make the commit order appear out of order ([#4696](https://github.com/git-town/git-town/issues/4696)).

#### Contributors

Shoutout to @erik-rw, @kevgo, @legeana, @nekitk, @pradeepmurugesan, @ruudk, @stephenwade, @terheyden for contributing code, ideas, and feedback to 33 shipped PRs and 10 resolved issues!

## 18.3.2 (2025-04-09)

#### Bug Fixes

- Restores the previous behavior of attempting an initial push when syncing tracking branches ([#4681](https://github.com/git-town/git-town/issues/4681)).

#### Contributors

Shoutout to @kevgo, @ruudk for contributing bug reports, ideas, and code to 12 shipped PRs and 2 resolved issues!

## 18.3.1 (2025-04-07)

#### Bug Fixes

- All `git rebase` commands now consistently use the `--no-update-refs` flag ([#4678](https://github.com/git-town/git-town/issues/4678)).

#### Contributors

Shoutout to @kevgo, @ruudk for contributing feedback and code to 11 shipped PRs and 1 resolved issues!

## 18.3.0 (2025-04-05)

#### New Features

- The new [git town detach](https://www.git-town.com/commands/detach.html) command removes a branch from its stack and makes it an independent top-level branch. This allows you to review and ship more of your branches concurrently, and focuses stacks on changes that belong together ([#4620](https://github.com/git-town/git-town/pull/4620)).
- The new [git town swap](https://www.git-town.com/commands/swap.html) command switches the position of the current branch with its parent, i.e. moves the current branch one position forward in the stack. This allows you to group related branches together, for example to ship them together or [merge](https://www.git-town.com/commands/swap.md) them.
- [git town merge](https://www.git-town.com/commands/merge.html) no longer syncs branches on its own, and now requires all affected branches to be in sync. This separates merge conflicts arising from syncing from merge conflicts arising from merging. `git town merge` now effectively only deletes the parent branch ([#4655](https://github.com/git-town/git-town/pull/4655)).
- The help screen printed by Git Town commands now gives a usage example ([#4672](https://github.com/git-town/git-town/pull/4672)).

#### Bug Fixes

- `git town set-parent` no longer accidentally deletes commits of the branch in certain edge cases ([#4669](https://github.com/git-town/git-town/issues/4669)).
- `git town set-parent` no longer deletes conflicting files ([#4638](https://github.com/git-town/git-town/issues/4638)).
- `git town merge` now errors if the parent branch has more than one child ([#4658](https://github.com/git-town/git-town/pull/4658)).
- `git town prepend --beam` now prints the correct branch name in the picker dialog ([#4642](https://github.com/git-town/git-town/issues/4642)).
- `git town compress` now handles merge commits correctly ([#4563](https://github.com/git-town/git-town/issues/4563)).

#### Contributors

Shoutout to @FirelightFlagboy, @avaz, @charlierudolph, @cjol, @davidolrik, @erik-rw, @ianjsikes, @kevgo, @leonhfr, @levrik, @ruudk, @stephenwade, @tranhl for contributing ideas, feedback, and code to 55 shipped PRs and 10 resolved issues!

## 18.2.0 (2025-03-28)

#### New Features

- Git Town now lets you submit staged changes into a new branch using a single command. The [hack](https://www.git-town.com/commands/hack.html), [append](https://www.git-town.com/commands/append.html), and [prepend](https://www.git-town.com/commands/prepend.html) commands now have a `-c`/`--commit` flag to commit the staged changes into the new branch. Use `-m`/`--message` to specify a commit message. The new `--propose` flag goes one step further and immediately proposes the new branch. Both `--message` and `--propose` imply `--commit`, so you can now run `git hack bugfix --propose` to commit your staged changes into a new `bugfix` branch and create a pull request for it in one step ([#4376](https://github.com/git-town/git-town/issues/4376)).
- [git town sync](https://www.git-town.com/commands/sync.html) now has a [--prune](https://www.git-town.com/commands/sync.html#-p--prune) flag that removes branches with no changes.
- Git Town now works with repositories hosted on [Codeberg](https://codeberg.org).

#### Contributors

Shoutout to @WhosNickDoglio, @andrew-rosca, @avaz, @caccavale, @charlierudolph, @kevgo, @lud-wj, @ruudk, @stefanfrede, @stephenwade for contributing ideas, feedback, and code to 25 shipped PRs and 5 resolved issues!

## 18.1.0 (2025-03-20)

#### New Features

- This release replaces the term "hosting platform" with [forge](https://en.wikipedia.org/wiki/Forge_(software)) because the latter is more correct and only one word. The configuration setting `hosting.platform` is now `hosting.forge-type`. This isn't a breaking change since the old settings still work. ([#4565](https://github.com/git-town/git-town/pull/4565), [#4568](https://github.com/git-town/git-town/pull/4568), [#4570](https://github.com/git-town/git-town/pull/4570))

#### Bug Fixes

- Pushing tags now also follows the [push-hook](https://www.git-town.com/preferences/push-hook.html) setting ([#4584](https://github.com/git-town/git-town/issues/4584)).
- Git Town now handles symbolic refs correctly ([#4588](https://github.com/git-town/git-town/issues/4588)).
- [git town undo](https://www.git-town.com/commands/undo.html) now unstashes at most one Git stash ([#4577](https://github.com/git-town/git-town/issues/4577)).

#### Contributors

Shoutout to @andrei9669, @blarson-hearst, @caccavale, @kevgo, @lud-wj for contributing ideas, feedback, and code to 24 shipped PRs and 4 resolved issues!

## 18.0.0 (2025-01-27)

#### BREAKING CHANGES

- `git town sync`: Local changes in a [prototype branch](https://www.git-town.com/branch-types#prototype-branches) now get pushed if that branch has a tracking branch. `git town sync` doesn't create this tracking branch, you have to create it manually if this behavior is needed. Prototoype branches are used for active development, hence their local and tracking branch should remain in sync ([#4542](https://github.com/git-town/git-town/issues/4542)).

#### Bug Fixes

- `git town compress` no longer removes the commit message body ([#4536](https://github.com/git-town/git-town/issues/4536)).
- No longer tries to look up proposals when [offline mode](https://www.git-town.com/preferences/offline.html) is enabled ([#4544](https://github.com/git-town/git-town/issues/4544)).
- `git town switch` now separates warning messages with an empty line from the branch tree ([#4543](https://github.com/git-town/git-town/issues/4543)).
- The [setup assistant](https://www.git-town.com/commands/config-setup.html) now allows configuring the new [ff-only sync strategy](https://www.git-town.com/preferences/sync-perennial-strategy.html#ff-only) ([#4549](https://github.com/git-town/git-town/pull/4549)).

#### Contributors

Shoutout to @JaredSharplin, @kevgo, @stephenwade for contributing ideas, feedback, and code to 17 shipped PRs and 4 resolved issues!

## 17.3.0 (2025-01-27)

#### New Features

- The new [ff-only sync strategy](https://www.git-town.com/preferences/sync-perennial-strategy.html#ff-only) provides a more elegant way to sync perennial branches that are protected against pushes and therefore should not receive local commits. Git Town fast-forwards the local branch to match the tracking branch. If a fast-forward is not possible, Git Town exits with a descriptive error message. This is ideal when you want an explicit warning about unpushed local commits ([#4104](https://github.com/git-town/git-town/issues/4104)).
- Git Town commands now uses the command listed in the `BROWSER` environment variable to open a browser. If no such environment variable exists, it opens the browser as before ([#4495](https://github.com/git-town/git-town/pull/4495)).
- git sync: syncs perennial, contribution, and observed branches without changes faster ([#4510](https://github.com/git-town/git-town/pull/4510), [#4513](https://github.com/git-town/git-town/pull/4513)).

#### Bug Fixes

- git sync: `--stack` now syncs observed branches correctly ([#4518](https://github.com/git-town/git-town/issues/4518)).

#### Contributors

Shoutout to @davidolrik, @FirelightFlagboy, @kevgo, @legeana, @sergej-koscejev, @stephenwade, @tugrulates, @wayne-zhan, @wlohrmann-hf for contributing ideas, feedback, and code to 56 shipped PRs and 10 resolved issues!

## 17.2.0 (2025-01-02)

#### New Features

- The new [always-merge ship strategy](https://www.git-town.com/preferences/ship-strategy.html#always-merge) always creates a merge commit when shipping a branch ([#4381](https://github.com/git-town/git-town/issues/4381)).
- `git town prepend` now has a [--beam](https://www.git-town.com/commands/prepend.html#-b--beam) option. When enabled, it allows the user to choose one or more commits to move into the new branch that is getting prepended ([#4356](https://github.com/git-town/git-town/pull/4356)).
- `git town prepend` now has a [--propose](https://www.git-town.com/commands/prepend.html#--propose) flag. When enabled, it automatically proposes the new branch. This is meant to be used together with `--beam` ([#4377](https://github.com/git-town/git-town/issues/4377)).
- Git Town's `hack` command can now make any branch type a feature branch ([#4373](https://github.com/git-town/git-town/issues/4373)).
- The new [git town status show command](https://www.git-town.com/commands/status-show.html) displays details about the currently interrupted Git Town command ([#4457](https://github.com/git-town/git-town/pull/4457)).
- Git Town now adds a message to entries it creates on the Git stash. This helps identify which stash entries were created by Git Town ([#4479](https://github.com/git-town/git-town/pull/4479)).
- If the only problem Git Town encounters is conflicts when running `git stash pop` at the end, it drops the stash entry and no longer returns with an error ([#4480](https://github.com/git-town/git-town/pull/4480)).

#### Bug Fixes

- The setup assistant no longer crashes when run in a brand-new repository ([#4365](https://github.com/git-town/git-town/pull/4365)).
- `git town status reset` now works ([#4469](https://github.com/git-town/git-town/pull/4469)).

#### Contributors

Shoutout to @Crocmagnon, @kevgo, @legeana, @lud-wj, @ruudk, @stephenwade, @wayne-zhan for contributing feedback, ideas, documentation, and code to 104 shipped PRs and 7 resolved issues!

## 17.1.1 (2024-12-20)

Git Town 17.1.1 ships a bug fix that unblocks BitBucket Datacenter users.

#### Bug Fixes

- Fixes a bug in the BitBucket-Datacenter driver ([#4371](https://github.com/git-town/git-town/pull/4371)).

#### Contributors

Shoutout to @Crocmagnon, @kevgo, @stephenwade for contributing bug fixes and ideas to 13 shipped PRs.

## 17.1.0 (2024-12-19)

Git Town 17.1 unblocks users who have submitted urgent issues.

#### New Features

- The configuration file can now also be named ".git-town.toml", in addition to ".git-branches.toml". This helps users locate it when looking for "Git Town configuration". ([#4343](https://github.com/git-town/git-town/issues/4343)).
- Supports development versions of Git ([#4344](https://github.com/git-town/git-town/pull/4344)).
- `git town switch` no longer asks for ancestry information. This avoids the risk of accidentally selecting a wrong parent branch when one is under the assumption of selecting a branch to switch to ([#4340](https://github.com/git-town/git-town/issues/4340)).
- Git Town now uses full-length SHA values to address commits. This improves reliability in very large repositories ([#4339](https://github.com/git-town/git-town/pull/4339)).

#### Bug Fixes

- The setup assistant now allows configuring a "BitBucket-Datacenter" setup ([#4360](https://github.com/git-town/git-town/pull/4360)).

#### Contributors

Shoutout to @kevgo, @lud-wj, @pratikpc, @ruudk, @stephenwade for contributing feedback, ideas, documentation, and code to 28 shipped PRs and 5 resolved issues!

## 17.0.0 (2024-12-16)

Git Town v17 modernizes some outdated concepts that were already documented. As always with major new Git Town releases, all changes are backwards compatible, so Git Town users don't need to change anything. The major version bump is merely a signal that you might need to update some of your own tooling around Git Town.

#### BREAKING CHANGES

- The configuration setting `create-prototype-branches` gets generalized into `new-branch-type`, which allows giving new branches all available branch types ([#3913](https://github.com/git-town/git-town/issues/3913)).
- The output of `git town config` now displays data organized in more sections ([#3866](https://github.com/git-town/git-town/issues/3866)).
- The config file gets generated without comments now, improving its readability ([#4335](https://github.com/git-town/git-town/pull/4335))
- The output of `git town config` now indicates more precisely whether an configuration setting is not provided or empty ([#4265](https://github.com/git-town/git-town/pull/4265)).

#### New Features

- When removing a branch, `git town sync` and `git town delete` now also remove the changes of those branches from their descendents ([#4189](https://github.com/git-town/git-town/issues/4189)).
- `git town set-parent` now also removes the changes from former parent branches ([#3473](https://github.com/git-town/git-town/issues/3473)).
- Git Town can now auto-resolve merge conflicts that include deleted files ([#4289](https://github.com/git-town/git-town/pull/4289)).
- The name of the Git remote used for development is now configurable ([#4165](https://github.com/git-town/git-town/issues/4165)).
- The setup assistant (`git town config setup`) now asks for the [sync-prototype-strategy](https://www.git-town.com/preferences/sync-prototype-strategy).
- `git town config` now displays the [sync-prototype-strategy](https://www.git-town.com/preferences/sync-prototype-strategy).

#### Bug Fixes

- Git Town no longer pops up the editor when continuing a command that got stuck in a rebase in certain situations ([#4285](https://github.com/git-town/git-town/pull/4285)).
- Now only updates Git aliases when auto-updating of outdated configuration entries that aren't Git Town settings ([#4304](https://github.com/git-town/git-town/pull/4304)).

#### Contributors

Shoutout to @Crocmagnon, @WhosNickDoglio, @alphatroya, @avaz, @erik-rw, @kevgo, @koppor, @matthewgonzalez, @mklauer, @nishchay-manwani-glean, @pandvan, @ruudk, @sheldonhull, @smaldored, @stephenwade for contributing feedback, ideas, documentation, and code to 70 shipped PRs and 13 resolved issues!

## 16.7.0 (2024-11-30)

#### New Features

- Support for BitBucket Datacenter ([#4239](https://github.com/git-town/git-town/pull/4239))

#### Bug Fixes

- Fixes a regression where branches shipped at the remote were no longer cleaned up ([#4222](https://github.com/git-town/git-town/issues/4222))
- Fixes a regression where branches that were shipped at the remote and deleted manually on the local machine were no longer cleaned up ([#4132](https://github.com/git-town/git-town/issues/4132))
- Correctly handles leading and trailing spaces in Git metadata branch lists ([#4240](https://github.com/git-town/git-town/pull/4240)
- `git town status reset` now works even if the persisted status is invalid ([#4236](https://github.com/git-town/git-town/pull/4236))

#### Contributors

Shoutout to @Crocmagnon, @alexdavid, @kevgo, @ruudk, @stephenwade for contributing feedback, ideas, documentation, and code to 26 shipped pull requests and 6 resolved issues!

## 16.6.1 (2024-11-20)

#### Bug Fixes

- `git town ship` no longer closes child proposals when using the `fast-forward` ship-strategy ([#4210](https://github.com/git-town/git-town/issues/4210))
- `git town propose` no longer opens the browser when proposing a branch that was deleted at the remote ([#4171](https://github.com/git-town/git-town/issues/4171))
- Bash-based installer makes the binary executable ([#4231](https://github.com/git-town/git-town/pull/4231))
- Stability improvements ([#4216](https://github.com/git-town/git-town/pull/4216))

#### Contributors

Shoutout to @EngHabu, @kevgo, @stephenwade, @thatch, @tranhl for contributing feedback, ideas, documentation, and code to 20 shipped PRs and 6 resolved issues!

## 16.6.0 (2024-11-12)

#### New Features

- Git Town now auto-resolves phantom merge conflicts for the `merge` and `compress` [sync-feature-strategy](https://www.git-town.com/preferences/sync-feature-strategy) ([#4183](https://github.com/git-town/git-town/pull/4183)). This eliminates the need to manually resolve unnecessary merge conflicts when you ship the oldest branch in a stack and then sync the rest of the stack.
- The new [git town merge](https://www.git-town.com/commands/merge) command merges two adjacent branches in a branch stack ([#4196](https://github.com/git-town/git-town/pull/4196)).

#### Bug Fixes

- Git Town's new capabilities-driven architecture now only attempts to change proposals if the hosting platform supports it. This reduces error messages by preventing unsupported operations ([#4203](https://github.com/git-town/git-town/pull/4203), [#4200](https://github.com/git-town/git-town/pull/4200)).
- Git Town now always syncs a branch with its parent branch before syncing with its tracking branch ([#4193](https://github.com/git-town/git-town/pull/4193)).

#### Contributors

Big thanks to @IGassmann, @ceilfors, @heisencoder, @kevgo, @mball-agathos, @stephenwade, @tranhl for contributing invaluable ideas to 30 shipped PRs and 2 resolved issues.

## 16.5.0 (2024-10-25)

#### New Features

- `git town sync` now automatically deletes a local branch if its remote tracking branch is gone, without attempting to sync it first--assuming the branch was in sync before Git Town ran. ([#3641](https://github.com/git-town/git-town/issues/3641)).
- The new `git town rename` command replaces the current `git town rename-branch` command. The `rename-branch` command is now deprecated and will be removed in a future release. Please update your tooling to use `git town rename` from now on. Existing Git aliases calling `git town rename-branch` get automatically updated to the new commands ([#4038](https://github.com/git-town/git-town/issues/4038)).
- Same for the new `delete` command, it replaces the now deprecated `kill` command. ([#4039](https://github.com/git-town/git-town/issues/4039)).
- You can now run `git town branch` in the middle of an unfinished sync ([#4108](https://github.com/git-town/git-town/issues/4108)).

#### Bug Fixes

- `git town ship` is now resilient against concurrently running Git processes ([#4142](https://github.com/git-town/git-town/pull/4142)).
- `git town propose` now pushes prototype branches after converting them to feature branches ([#4109](https://github.com/git-town/git-town/issues/4109)).
- `git town propose` now uses the first existing ancestor branch as the parent if the current parent branch was shipped or deleted remotely ([#4135](https://github.com/git-town/git-town/pull/4135)).

#### Contributors

Big thanks to @JamieMcKernanKaizen, @PowerSchill, @bengeois, @kevgo, @ruudk, @sergej-koscejev, @stephenwade, @tranhl, @vectro for contributing to 53 shipped PRs and 12 resolved issues!

## 16.4.1 (2024-10-09)

Thanks to @bengeois, @kevgo, @ruudk, @sergej-koscejev, @stephenwade, @tranhl for contributing great feedback, ideas, and code to 14 shipped PRs and 5 resolved issues!

#### Bug Fixes

- disable `rebase.updateRefs` if your Git version supports it ([#4101](https://github.com/git-town/git-town/discussions/4101))

## 16.4.0 (2024-10-03)

Git Town 16.4 improves the usability in busy monorepos as well as the stability and resilience of Git Town in more environments.

Many thanks to @FirelightFlagboy, @JamieMcKernanKaizen, @alexw10, @blaggacao, @charlierudolph, @ericcao52, @ianjsikes, @kevgo, @ruudk, @sheldonhull, @shiv19, @stephenwade, @tranhl, @waldyrious for improving Git Town through 39 shipped PRs and 13 resolved issues!

#### New Features

- Git Town's `append`, `hack`, `prepend`, and `propose` commands now have a `--detached` flag that prevents them from pulling in additional changes from the main branch. This together with the existing `--detached` flag for `git town sync` allows controlling exactly when changes from the main branch get synced into your branch hierarchy ([#4095](https://github.com/git-town/git-town/issues/4059)).
- New config settings [contribution-regex](https://www.git-town.com/preferences/contribution-regex) and [observed-regex](https://www.git-town.com/preferences/observed-regex) allow marking branches created by external services like Renovate or Dependabot appropriately ([#3985](https://github.com/git-town/git-town/issues/3985)).
- Git Town is now much more resilient against unexpected Git failures, for example when another Git process is running concurrently, because most Git Town operations are now fully reentrant ([#4082](https://github.com/git-town/git-town/pull/4082)).
- `git town sync` now syncs branches whose ancestors aren't available locally better: It pulls the tracking branches of all non-local ancestors until it finds a local ancestor ([#3769](https://github.com/git-town/git-town/issues/3769)).

#### Bug Fixes

- `git town hack` no longer panics if the main branch isn't available locally ([#3703](https://github.com/git-town/git-town/issues/3703)).
- `git town hack` no longer panics when the Git repo has a detached HEAD ([#3694](https://github.com/git-town/git-town/issues/3694)).
- Git Town now loads all applicable settings from the config file ([#4072](https://github.com/git-town/git-town/issues/4072)).

## 16.3.0 (2024-09-24)

Git Town 16.3 is here, and it's packed with some long-requested features we've been working towards for years!

Big shoutout to @LogvinovLeon, @alexw10, @charlierudolph, @cjol, @ericcao52, @kevgo, @mw00120, @ruudk, @stephenwade, @tranhl, @waldyrious, @zodman for your insightful contributions, feedback, and ideas. Git Town is a much better product thanks to you and the 52 shipped PRs and 6 resolved issues in this release!

#### New Features

- Full integration of proposals ([#2745](https://github.com/git-town/git-town/issues/2745)):
  - Git Town now updates all affected pull requests when you rename, prepend, or remove a branch or change its parent.
  - If the parent of a branch is unknown but there's an open PR, Git Town will now grab the PR's base branch as the parent.
  - `git town undo` now also reverts any changes made to pull requests ([#4049](https://github.com/git-town/git-town/issues/4049)).
- Full integration with the Bitbucket Cloud API ([#973](https://github.com/git-town/git-town/issues/973)) and the gitea API ([#4044](https://github.com/git-town/git-town/pull/4044)).
- `git town rename-branch` now maintains the Git configuration and reflog for renamed branches ([#4023](https://github.com/git-town/git-town/issues/4023)).
- Git Town now logs all API interactions in the CLI, showing details like PR numbers and branch names it retrieves from the APIs. This makes it easier to see where information and possible slowness come from ([#4020](https://github.com/git-town/git-town/pull/4020), [#4026](https://github.com/git-town/git-town/pull/4026)).

## 16.2.1 (2024-09-12)

Thanks to @kevgo, @ruudk, @stephenwade, @zodman for reporting and fixing unintuitive behavior that can and should be improved, leading to 13 shipped PRs and 2 resolved issues!

#### Bug Fixes

- `git town switch` no longer asks for the ancestry of unrelated branches ([#4004](https://github.com/git-town/git-town/issues/4004)).
- `git town branch` no longer prints a redundant newline ([#4011](https://github.com/git-town/git-town/pull/4011)).

## 16.2.0 (2024-09-12)

Git Town 16.2 makes it easier for you to manage the branches in your local repository.

Big thanks to @kevgo, @ruudk, @stephenwade, @zodman for coming up with the ideas for these new features and helping polish them in 38 shipped PRs and 6 resolved issues!

#### New Features

- The new [git town branch](https://www.git-town.com/commands/branch) command displays the local branch hierarchy, and the types of all branches except for main and feature branches ([#3807](https://github.com/git-town/git-town/issues/3807)).
- `git town switch` can now displays the types of branches (except for the main and feature branches) when called with the [--display-types](https://www.git-town.com/commands/switch#--display-types---d) flag ([#3937](https://github.com/git-town/git-town/issues/3937)).
- `git town switch` can now filter the branches to switch to via [regular expressions](https://www.git-town.com/commands/switch#positional-arguments) ((#3980)[https://github.com/git-town/git-town/pull/3980/files])
- You can now use `git town switch` to check out a remote branch using the [--all](https://www.git-town.com/commands/switch#--all---a) flag ([#3941](https://github.com/git-town/git-town/issues/3941))

#### Bug Fixes

- `git town switch` now exits with a nice error message if there are no branches to switch to ([#3979](https://github.com/git-town/git-town/issues/3979))
- logs of API calls are now capitalized ([#3975](https://github.com/git-town/git-town/issues/3975))

## 16.1.1 (2024-09-09)

Git Town 16.1.1 ships important bug fixes for the new features introduced in v16.1.

Thanks to @kevgo, @rbardini, @stephenwade!

#### Bug Fixes

- The setup assistant now always stores `default-branch-type` and `feature-regex` in the Git-based configuration, not in the config file. These settings are typically developer-specific. You can still manually add them to the config file if needed ([#3961](https://github.com/git-town/git-town/pull/3961)).
- Fixes various documentation errors ([#3953](https://github.com/git-town/git-town/pull/3953)).

## 16.1.0 (2024-09-07)

Git Town 16.1 provides multiple killer features for power users.

Big thanks to @FFdhorkin, @breml, @bryanlarsen, @buscape, @enigma, @heyitsaamir, @kevgo, @rbardini, @ruudk, @stephenwade, @tranhl, @zeronacer for contributing really good ideas, feedback, and code to 39 shipped PRs and 9 resolved issues!

#### New Features

- **Automatic branch detection:** Git Town can now automatically detect feature and contribution or observed branches if you configure the new [default-branch-type](https://www.git-town.com/preferences/default-branch-type) and [feature-regex](https://www.git-town.com/preferences/feature-regex) options ([#3683](https://github.com/git-town/git-town/issues/3683)).
- **Detached syncing:** Use `git sync --detached` to sync without pulling updates from the main branch. This helps keep development momentum if the main branch receives frequent updates and these updates trigger costly follow-up activities like `npm install` or long-running recompiles ([more info](https://www.git-town.com/commands/sync), [#2657](https://github.com/git-town/git-town/issues/2657)).
- **More concise branch switching:** Running `git town switch --type=<branch types>` displays only branches of the given type ([more info](https://www.git-town.com/commands/switch), [#3933](https://github.com/git-town/git-town/issues/3933)).

#### Bug Fixes

- Git Town no longer errors if the previously checked out Git branch is active in another worktree ([#3916](https://github.com/git-town/git-town/pull/3916)).

## 16.0.0 (2024-08-30)

Git Town 16 brings the "git ship" command back in a big way: Git Town now supports shipping stacked changes without merge conflicts - even on platforms that don't support it natively, like GitHub!

Big thanks to @FFdhorkin, @antoineMoPa, @breml, @bryanlarsen, @buscape, @kevgo, @tranhl, @zeronacer for the great feedback that led to this awesome new solution! This releaese contains 9 shipped PRs and 7 resolved issues.

#### BREAKING CHANGES

The default behavior of `git ship` tightens. Previously it shipped via the API if an API key is configured, and without an API key it did a local squash-merge. The new default behavior is to ship only via API or not at all. The new default behavior is safer because it only automates what the user would normally do online. You can specify a different behavior for `git ship` via the new `ship-strategy` configuration option (see below).

#### New Features

- You can now configure how Git Town ships branches via the new `ship-strategy` configuration setting. Possible options are:
  - `api` ships the branch by merging its proposal via the API of your code hosting platform.
  - `fast-forward` is a new shipping strategy that prevents the false merge conflicts you get when shipping a branch from a stack using squashes or merges. It merges the branch to ship via `git merge --ff-only` into its parent (typically the main branch) on your local machine and then pushes the new commits to the remote main branch.
  - `squash-merge` as before merges the branch to ship via `git merge --squash` into its parent.

## 15.3.0 (2024-08-26)

Git Town 15.3 brings sweet quality-of-life improvements.

Massive thanks to @ChiefMateStarbuck, @IvanVas, @WhosNickDoglio, @alphatroya, @charlierudolph, @cirego, @erik-rw, @gstamp, @guusw, @kelbyers, @kevgo, @marcelpanse, @nishchay-manwani-glean, @rnystrom, @ruigomeseu, @sergej-koscejev for helping evolve Git Town with useful feedback, ideas, and code contributions to 21 shipped PRs and 14 resolved issues!

#### New Features

- Automatic retry for concurrent Git access: Git Town now waits and retries Git operations if another Git process is running concurrently. Super handy when your IDEs is running Git commands in the background ([#3629](https://github.com/git-town/git-town/issues/3629)).
- Shell prompt status indicator: If a Git Town command gets interrupted by a merge conflict, you can now add the name of the pending Git Town command to your shell prompt. This reminds you to run `git town continue` to finish the job ([#2208](https://github.com/git-town/git-town/issues/2208)).
- `git town propose` now takes you directly to the existing proposal's webpage if one already exists ([#2362](https://github.com/git-town/git-town/issues/2362)).
- API activity logs. Git Town now logs its communication with hosting APIs in the CLI output. This shows you what Git Town is doing and where slowness is coming from ([#3892](https://github.com/git-town/git-town/pull/3892)).
- `git town kill` no longer asks for the ancestry of branches it is about to delete as long as these branches don't have descendents ([#3870](https://github.com/git-town/git-town/issues/3870)).
- Setting up shell autocompletion on ZSH is now better documented ([#3889](https://github.com/git-town/git-town/pull/3889)).

#### Bug Fixes

- If you have a merge conflict between your uncommitted changes and branch ancestry, Git Town commands will fail when running `git stash pop` at the end. Previously when running `git town continue` it tried to pop the stash again, causing the same merge conflict to happen again. Now Git Town assumes you have resolved the merge conflicts when running `git town continue` and deletes the stash entry. If you need to re-apply the conflicting stash entry, run `git stash pop` manually before running `git town continue`. This keeps your Git stash clean ([#3886](https://github.com/git-town/git-town/pull/3886)).
- `git town continue` now re-runs all failed Git operations, helping recover from a wider range of unexpected issues ([#3887](https://github.com/git-town/git-town/pull/3887), [#3885](https://github.com/git-town/git-town/pull/3885)).

## 15.2.0 (2024-08-21)

Big thanks to @kevgo, @mball-agathos, @ruudk, @sergej-koscejev for contributing super useful feedback, ideas, and code to 32 shipped PRs and 5 resolved issues!

#### New Features

- The new "compress" sync strategy always compresses branches while syncing them ([#3320](https://github.com/git-town/git-town/issues/3320)).
- Basic support for integrating Git Town into [lazygit](https://github.com/jesseduffield/lazygit) ([#3872](https://github.com/git-town/git-town/pull/3872)).

#### Bug Fixes

- Renaming branches now keeps their contribution, observed, parked, and prototype status ([#3864](https://github.com/git-town/git-town/issues/3864)).
- The commands git town contribute, observe, park, and prototype now behave more correct and consistent ([#3880](https://github.com/git-town/git-town/pull/3880)).

## 15.1.0 (2024-08-09)

Numerous thanks to @FirelightFlagboy, @Iron-Ham, @IvanVas, @JaredSharplin, @JustinBis, @TheHolyWaffle, @WhosNickDoglio, @alexus37, @alphatroya, @anikrajc, @blaggacao, @charlierudolph, @cjol, @connected-rmcleod, @cridasilva, @defunctzombie, @erik-rw, @kevgo, @pattiereaves, @sgarfinkel, @stephenwade, @teumas, @zodman for the super useful feedback, ideas, and code contributions to 31 shipped PRs and 19 resolved issues.

#### New Features

- `git repo` can now take the name of a remote to open the repo at that remote ([#1204](https://github.com/git-town/git-town/issues/1204)).
- The new `sync-tags` config option disables syncing of Git tags ([#3212](https://github.com/git-town/git-town/issues/3212)).
- `git ship` can now ship into any type of parent branch with the `--to-parent` option ([#2605](https://github.com/git-town/git-town/issues/2605)).
- `git sync --stack` syncs all branches in the current stack ([#3816](https://github.com/git-town/git-town/pull/3816)).

## 15.0.0 (2024-08-05)

Git Town 15.0 improves Git Town's compatibility with monorepos and removes technical drift.

Major thanks to @ianjsikes, @kevgo, @ruudk, @seadowg, @stephenwade, @zodman for contributing valuable feedback, ideas, and code to 41 shipped PRs and 8 resolved issues!

#### BREAKING CHANGES

- `git town ship` no longer syncs branches when shipping. From now on it only ships branches that are in sync. This ensures that only fully tested and reviewed changes get shipped ([#3350](https://github.com/git-town/git-town/issues/3350)).
- This also makes the `sync-before-ship` config option obsolete, it no longer exists ([#3644](https://github.com/git-town/git-town/pull/3644)).
- `git town prepend` no longer syncs when uncommitted changes are present. This allows committing your uncommitted changes first, then syncing later ([#3778](https://github.com/git-town/git-town/pull/3778)).
- The term `main development branch` gets shortened to `main branch` since there are no other development branches in Git Town's domain model ([#3643](https://github.com/git-town/git-town/issues/3643)).

#### New Features

- A new branch type called _prototype branches_ syncs only locally, i.e. they don't create or push to a tracking branch until they are proposed. This helps reduce stress on the CI server, allows developers to prototype using sensitive information or potentially problematic code or data that they don't want to share ([#3646](https://github.com/git-town/git-town/issues/3646)).
- The new `sync-prototype-strategy` setting allows defining a dedicated sync strategy for prototype branches. This allows rebasing your commits while they are local, and switching to merging once other developers can see them ([#3785](https://github.com/git-town/git-town/pull/3785)).
- The new `create-prototype-branches` setting makes Git Town always create prototype branches ([#3779](https://github.com/git-town/git-town/pull/3779)).

## 14.4.1 (2024-07-29)

Many thanks to @charlierudolph, @ianjsikes, @kevgo, @seadowg, @stephenwade for contributing feedback, ideas, and code to 10 shipped PRs and 4 resolved issues.

### Bug Fixes

- `git town undo` now only undoes changes to branches that the previous Git Town command has touched ([#3765](https://github.com/git-town/git-town/issues/3765))
- `git town continue` now does not ask for additional lineage information ([#3725](https://github.com/git-town/git-town/issues/3725))

## 14.4.0 (2024-07-26)

Git Town v14.4 ships features suggested by Git Town users and fixes a severe bug in the undo feature.

Many thanks to @ianjsikes, @kevgo, @ruudk, @stephenwade who contributed ideas for great new features and helped identify and solve a severe bug, resulting in 46 shipped PRs and 2 resolved issues.

#### New Features

- Prototype branches are fully synced branches that don't have a tracking branch. They are useful when working with sensitive information or to save on CI minutes ([#3646](https://github.com/git-town/git-town/issues/3646)).

#### Bug Fixes

- Undo only local branches and their tracking branches ([#3764](https://github.com/git-town/git-town/pull/3764)).
- Use correct parent when deleting shipped observed branches ([#3757](https://github.com/git-town/git-town/pull/3757), [#3756](https://github.com/git-town/git-town/pull/3756), [#3755](https://github.com/git-town/git-town/pull/3755), [#3754](https://github.com/git-town/git-town/pull/3754/files)).
- Undo for configuration commands now correctly undoes branch changes ([#3741](https://github.com/git-town/git-town/pull/3741/files)).

## 14.3.1 (2024-07-15)

#### Bug Fixes

- `git sync --no-push` now also doesn't push when the `rebase` sync-strategy is configured ([#3271](https://github.com/git-town/git-town/pull/3721))
- `git town config get-parent` no longer prints an unnecessary empty line ([#3717](https://github.com/git-town/git-town/pull/3717))

Many thanks to @dannykingme, @defunctzombie, @kevgo, @marcosfelt, @nekitk, @opeik, @pcfreak30, @ruudk, and @stephenwade for identifying the issues fixed in this release and providing helpful feedback to resolve them, resulting in 19 shipped PRs and 4 resolved issues!

## 14.3.0 (2024-07-12)

This release ships a few of the most requested community features.

#### New Features

- `git propose` now supports flags to pre-populate more fields of the pull requests to create:
  - `--title=<value>` sets the title to the given value
  - `--body=<value>` sets the body to the given value
  - `--body-file=<file path>` sets the body to the content of the file with the given path. Providing `-` as the file path reads the body from STDIN ([#3207](https://github.com/git-town/git-town/issues/3207)).
- a new command `git town config get-parent [branch]` prints the parent of the given branch, or the current branch if no branch is provided ([#3207](https://github.com/git-town/git-town/issues/3207)).
- the new `--no-push` flag for `git sync` temporarily disables pushing local changes when [sync-feature-strategy](https://www.git-town.com/preferences/sync-feature-strategy) is `merge`. Please note that when sync-feature-strategy is `rebase`, it still force-pushes to avoid keeping outdated commits around, which avoids data loss in edge cases.

Heartfelt thanks to @dannykingme, @defunctzombie, @kevgo, @marcosfelt, @nekitk, @opeik, @pcfreak30, @ruudk, @stephenwade for contributing code, ideas, and feedback to 29 shipped PRs and 9 resolved issues!

## 14.2.3 (2024-06-25)

Another release with bug fixes and stability improvements.

Shoutout to @alexus37, @bb010g, @blaggacao, @bryanlarsen, @charlierudolph, @kelbyers, @kevgo, @kinyat, @ruudk, @stephenwade, @vectro for contributing code, ideas, and feedback to 54 shipped PRs and 9 resolved issues!

#### Bug Fixes

- loads the configuration when calling Git Town from a subfolder in the Git repo ([#3688](https://github.com/git-town/git-town/pull/3688))
- "git observe" no longer panics when given a non-existing branch name ([#3647](https://github.com/git-town/git-town/issues/3647))

## 14.2.2 (2024-06-06)

This release fixes a few bugs that now get correctly surfaced thanks to stronger type checking introduced in v14.2.1.

Massive thanks to @breml, @bryanlarsen, @edwarbudiman, @FirelightFlagboy, @kelbyers, @kevgo for contributing code, ideas, and feedback to 13 shipped pull requests and 9 resolved issues!

#### Bug Fixes

- follows include directives in the Git configuration ([#3614](https://github.com/git-town/git-town/issues/3614))
- fixes a panic during git town propose ([#3539](https://github.com/git-town/git-town/issues/3539))

## 14.2.1 (2024-05-30)

This release brings substantial stability improvements due to much stronger type checking and removing unnecessary optionality and mutability from the codebase.

Big thanks to @breml, @bryanlarsen, @edwarbudiman, @FirelightFlagboy, @kevgo, @shiv19, @SophiaSaiada for contributing code, ideas, and feedback to 191 shipped pull requests and 9 resolved issues!

#### Bug Fixes

- fixes a panic when additional Git remotes are present ([#3537](https://github.com/git-town/git-town/issues/3537))
- fixes a panic when encountering invalid lineage entries ([#3453](https://github.com/git-town/git-town/issues/3453))
- fixes a panic if the previous Git branch is checked out in another worktree ([#3297](https://github.com/git-town/git-town/issues/3297))

## 14.2.0 (2024-04-24)

#### New Features

- `git town set-parent` is now a proper Git Town command, with continue and undo ([#3407](https://github.com/git-town/git-town/pull/3407)).
- Git Town now works if you have [merge.ff-only](https://git-scm.com/docs/git-merge#Documentation/git-merge.txt---ff-only) configured ([#3408](https://github.com/git-town/git-town/pull/3408)).

#### Bug Fixes

- `git town set-parent` now properly defaults to the existing parent ([#3406](https://github.com/git-town/git-town/pull/3406)).

Big thanks to @charlierudolph, @ericyliu, @hammenm, @hmbrg, @kevgo, @KORDayDream, @StevenXL, @vectro, @zifeo for contributing code, ideas, and feedback to 23 shipped PRs and 7 resolved issues!

## 14.1.0 (2024-04-19)

Besides polishing the `git town switch` command, Git Town 14.1 focuses on stability improvements and bashing bugs. We closed out 50% of all open tickets (70 tickets), including many long-standing bugs!

This version also de-emphasizes `git ship`. Most people should not run `git ship`. The recommended workflow is to ship feature branches using the web UI or merge queue of your code hosting platform. `git ship` is for edge cases like development in [offline mode](https://www.git-town.com/commands/offline).

#### New Features

- `git town switch` now has an `-m` option that checks out the selected branch using [git checkout -m](https://git-scm.com/docs/git-checkout#Documentation/git-checkout.txt--m) ([#3321](https://github.com/git-town/git-town/issues/3321)).
- `git town switch` now doesn't allow selecting branches that are checked out in other Git worktrees ([#3295](https://github.com/git-town/git-town/issues/3295)).
- `git town switch` now indicates the existence of uncommitted changes. This helps remember to commit them on the current branch if that was needed ([#3307](https://github.com/git-town/git-town/issues/3307)).
- Git Town now shuts down gracefully and allows continue and undo when you press `Ctrl-C` to cancel a Git command that runs too long or hangs ([#414](https://github.com/git-town/git-town/issues/414)).
- Notifications to the user are now highlighted in cyan in the Git Town output, making them easier to spot ([#3353](https://github.com/git-town/git-town/pull/3353)).
- The setup assistant now also uses `remotes/origin/HEAD` to determine the default main branch if the already used `init.defaultbranch` setting isn't set ([#646](https://github.com/git-town/git-town/issues/646)).
- Prototypical support for the API of GitHub Enterprise. This is impossible to test for the Git Town team, so please provide bug reports if something doesn't work ([#1179](https://github.com/git-town/git-town/issues/1179)).
- Improved support for GitLab instances that use a custom SSH port ([#1891](https://github.com/git-town/git-town/issues/1891)).
- `git ship` now sqash-merges using the `--ff` option. This removes an incompatibility for users who have the `merge.ff` option set to `false` in their Git configuration ([#1097](https://github.com/git-town/git-town/issues/1097)).
- If a branch is listed as its own parent, Git Town now notifies the user and deletes this invalid lineage entry ([#3393](https://github.com/git-town/git-town/pull/3393)).
- Improved error messages ([#2949](https://github.com/git-town/git-town/issues/2949)).

#### Bug Fixes

- `git sync` now ends on the previously checked out branch when pruning branches ([#2784](https://github.com/git-town/git-town/issues/2784)).
- `git sync --all` now syncs in topological order. This ensures all branches in deep stacks get synced ([#3344](https://github.com/git-town/git-town/pull/3344)).
- `git town switch` no longer displays branches that were deleted manually ([#3361](https://github.com/git-town/git-town/pull/3361)).
- `git kill` now checks out the main branch when the previous branch also was killed ([#3358](https://github.com/git-town/git-town/pull/3358)).

Massive thanks to @abhijeetbhagat, @aeneasr, @allewun, @alphatroya, @amarpatel, @avaz, @breml, @bryanlarsen, @charlierudolph, @ChrisMM, @cirego, @ericyliu, @grignaak, @hammenm, @hmbrg, @JCB-K, @kevgo, @koppor, @KORDayDream, @martinjaime, @mball-agathos, @mribichich, @ruudk, @sascha-andres, @sheldonhull, @tranhl, @vectro, @WhosNickDoglio, @WurmD, @zeronacer, @zifeo for contributing code, ideas, and feedback to 74 shipped PRs and 70 resolved issues!

## 14.0.0 (2024-04-12)

Git Town 14.0 improves the developer experience around uncommitted and stacked changes.

#### BREAKING CHANGES

`git hack`, `git append`, and `git prepend` no longer sync the branch lineage in the presence of uncommitted changes. This allows you to commit your changes first before pulling in more changes from other developers. They still sync if you call them without uncommitted changes ([#3198](https://github.com/git-town/git-town/issues/3198)).

#### New Features

`git town compress` (aliasable to `git compress` by re-running `git town config setup`) squashes all commits in a branch into a single commit. By default the new commit uses the commit message of the first commit in the branch. You can provide a custom commit message using the `-m` switch the same way as in `git commit`. The `--stack` option compresses all branches in a change stack. Git Town does not compress perennial, observed, contribution, and non-active parked branches ([#1529](https://github.com/git-town/git-town/issues/1529), [#2086](https://github.com/git-town/git-town/issues/2086)).

`git hack`, `git append`, and `git prepend` are faster due to creating and checking out the new branch using a single Git operation ([#3313](https://github.com/git-town/git-town/pull/3313)).

Big thanks to @blaggacao, @breml, @gabyx, @kevgo, @mball-agathos, @nishchay-manwani-glean, @pjh, @ruudk, @tranhl, @utkinn, and @WhosNickDoglio for contributing code, ideas, and feedback to 68 shipped PRs and 7 resolved issues!

## 13.0.2 (2024-03-29)

#### Bug Fixes

- Fixes a serious bug where users who have [branch.sort](https://git-scm.com/docs/git-branch#Documentation/git-branch.txt-branchsort) set might get tracking branches removed ([#3241](https://github.com/git-town/git-town/issues/3241)).

#### Statistics

Heartfelt thanks to @breml and @kevgo for going the extra mile to investigate a tricky bug and contributing to 25 shipped PRs and 4 resolved issues!

## 13.0.1 (2024-03-27)

#### Bug Fixes

- Allows syncing branches with merge conflicts in linked worktrees ([#3230](https://github.com/git-town/git-town/issues/3230)).
- Fixes the Bash-based installer ([#3234](https://github.com/git-town/git-town/pull/3234)).

#### Statistics

Many thanks to @alexus37, @breml, @bryanlarsen, @kevgo, @tranhl, @vectro, @wederbrand for contributing feedback, ideas, and solutions to 8 shipped PRs and 6 resolved issues!

## 13.0.0 (2024-03-22)

Git Town 13.0 adds better support for syncing feature branches after rebasing your commits and bumps the required Git version.

#### BREAKING CHANGES

When the [sync-feature-strategy](https://www.git-town.com/preferences/sync-feature-strategy) is set to `rebase`, Git Town now force-pushes your locally rebased commits to the tracking branch ([#3182](https://github.com/git-town/git-town/issues/3182)). This avoids mixing locally rebased commits with outdated commits on the tracking branch. To not accidentally override new commits on the tracking branch that haven't been integrated into your local commits, Git Town now force-pushes using the [--force-if-includes](https://git-scm.com/docs/git-push#Documentation/git-push.txt---no-force-if-includes) Git flag. This requires raising the minimally required Git version from 2.7 to 2.30. Git 2.30 was released over 2 years ago and should be widely available at this point.

#### New Features

- Git Town now automatically removes lineage entries for branches that were converted from feature branches to perennial branches ([#3218](https://github.com/git-town/git-town/issues/3218)).
- Git Town documentation and error messages now guide the user to call Git Town as `git town` instead of `git-town` on the CLI ([#3208](https://github.com/git-town/git-town/issues/3208)).

#### Bug Fixes

- Fixes a crash when an ustream HEAD is set ([#2660](https://github.com/git-town/git-town/issues/2660)).
- Fixes the error message when trying to set the parent of a perennial branch ([#3217](https://github.com/git-town/git-town/issues/3217)).

Kudos to @100rab-S, @dgentry, @kevgo, @koppor, @nicksieger, @ruudk, @srstevenson, and @tranhl for contributing code, ideas, and feedback to 18 shipped PRs and 13 resolved issues!

## 12.1.0 (2024-02-29)

Git Town 12.1 implements some of the most requested features by the Git Town community. It also continues the modernization of Git Town's internals. This time we made Git Town's undo engine simpler, more robust, and more reliable by removing all remaining mutable state.

#### New Features

- New options to fine-tune how Git Town syncs branches: `git contribute`, `git observe`, and `git park`. More info at https://www.git-town.com/advanced-syncing and [#3095](https://github.com/git-town/git-town/issues/3095).
- All branches matching the regular expression in the new configuration setting [branches.perennial-regex](https://www.git-town.com/preferences/perennial-regex) are now also considered perennial, in addition to the ones already listed in `branches.perennials`. This makes it easier to deal with situations where you have many perennial branches with similar sounding names like `release-1`, `release-2`, etc ([#2659](https://github.com/git-town/git-town/issues/2659)).

#### Bug Fixes

- [git town skip](https://www.git-town.com/commands/skip) now works correctly in complex situations ([#2978](https://github.com/git-town/git-town/issues/2978)).
- Git Town now deletes branches more reliably ([#3097](https://github.com/git-town/git-town/issues/3097)).

#### Statistics

Many thanks to @100rab-S, @harrismcc, @kevgo, @ruudk, @tranhl for contributing feedback, ideas, and solutions to 70 shipped PRs and 11 resolved issues!

## 12.0.2 (2024-02-14)

#### Bug Fixes

- All dialogs that show local branches now paginate ([#3119](https://github.com/git-town/git-town/issues/3119)).

## 12.0.1 (2024-02-12)

#### Bug Fixes

- removes crashes when using a self-hosted platform instance ([#3114](https://github.com/git-town/git-town/pull/3114))
- improve the CLI output when using a hosting connector ([#3115](https://github.com/git-town/git-town/pull/3115))

## 12.0.0 (2024-02-05)

Git Town 12 continues the effort to make the Git Town user experience more consistent and intuitive by modernizing Git Town's configuration system.

- 294 contributions
- 31 resolved tickets
- a heartfelt thanks to the contributors for this release: @alokpr, @brandonaut, @bryanlarsen, @ChrisMM, @eugef, @IGassmann, @Iron-Ham, @jakeleboeuf, @JaKXz, @kevgo, @koppor, @Nezteb, @ruudk, @zeronacer

#### BREAKING CHANGES

- The new setup assistent (see below) replaces the existing CLI and Git commands to change the configuration.
- Removing the Git Town configuration is now done by running `git town config remove` instead of `git town config reset` ([#3051](https://github.com/git-town/git-town/pull/3051)).
- More intuitive names for the following configuration options. Git Town automatically updates the configuration, so no action is needed on your end.
  - `code-hosting-platform` is now `hosting-platform` ([#3054](https://github.com/git-town/git-town/pull/3054))
  - `code-hosting-origin-hostname` is now `hosting-origin-hostname` ([#3053](https://github.com/git-town/git-town/pull/3053))
  - `ship-delete-remote-branch` is now `ship-delete-tracking-branch` ([#2929](https://github.com/git-town/git-town/pull/2929))
- Putting Git Town into offline mode is a top-level command again. Run `git town offline yes` to enable offline mode instead of `git town config offline yes` ([#3049](https://github.com/git-town/git-town/pull/3049)).
- All visual dialogs have been rewritten using a modern UI framework for a better look and to avoid the rendering issues encountered before ([#2964](https://github.com/git-town/git-town/issues/2964)).
- Nested feature branches are now called "stacked changes" to match the emerging industry term ([#3062](https://github.com/git-town/git-town/pull/3062)).

#### New Features

- Git Town v12 introduces Git Town's setup assistant ([#2941](https://github.com/git-town/git-town/issues/2941)). The setup assistant guides you through all of Git Town's configuration settings, including setting up the shorter aliases for Git Town commands. Run it by executing `git town config setup`. This assistant replaces the old configuration commands under `git town config`, the `alias` command, and the old "quick configuration" process.
- Git Town now supports storing non-confidential configuration entries in a configuration file with name `.git-branches.toml` ([#2748](https://github.com/git-town/git-town/issues/2748)). The best way to create one is the setup assistant. The setup assistant can also migrate your existing Git-based configuration to the config file.
- All commands now support the `--dry-run` flag to try them out safely ([#2859](https://github.com/git-town/git-town/pull/2859)).
- You can now install Git Town on Windows using Chocolatey: `choco install git-town` ([#763](https://github.com/git-town/git-town/issues/763))
- Massive performance improvements (exceeding 200%) on Windows thanks to not executing Git through the CMD shell anymore ([#2881](https://github.com/git-town/git-town/pull/2881)).
- The undo commands execute faster ([#2863](https://github.com/git-town/git-town/pull/2863)).

#### Bug Fixes

- Fix the `--version` command on Windows ([#2900](https://github.com/git-town/git-town/pull/2900/files)).

## 11.1.0 (2023-12-12)

#### New Features

- Git Town now handles branches checked out in other worktrees correctly ([#2764](https://github.com/git-town/git-town/pull/2764))
- Git Town now checks out the previous Git branch ("git checkout -") after removing a local branch ([#2742](https://github.com/git-town/git-town/pull/2742))

#### Bug Fixes

- `git continue` now correctly handles a manually popped stash after resolving conflicts ([#2758](https://github.com/git-town/git-town/pull/2758))
- `git continue` retries failing commit, merge-proposal, create-branch, create-proposal, create-tracking-branch, and push-branch operations ([#2756](https://github.com/git-town/git-town/pull/2756))
- `git continue` ensures there are no untracked files ([#2754](https://github.com/git-town/git-town/pull/2754))
- `git switch` now allows switching to perennial branches ([#2752](https://github.com/git-town/git-town/pull/2752))

## 11.0.0 (2023-12-06)

Git Town 11 continues the effort to make the Git Town user experience more consistent and intuitive.

#### BREAKING CHANGES

- `git new-pull-request` is now `git propose`. Not all platforms that Git Town supports use the name "pull request", so Git Town uses the word "proposal" for pull requests, merge requests, etc from now on. Nine fewer characters to type! ([#2691](https://github.com/git-town/git-town/pull/2691))
- `git abort` is merged into `git undo`. From now on you just run `git undo` after a Git Town command fails or finishes to get back to where you started ([#2719](https://github.com/git-town/git-town/pull/2719)).
- Many configuration options now have more intuitive names. No action needed on your end, Git Town automatically updates the affected settings on your machine. This means you can't go back to v10 after updating to v11.
  - `code-hosting-driver` is now `code-hosting-platform` ([#2704](https://github.com/git-town/git-town/pull/2704))
  - `main-branch-name` is now `main-branch` ([#2703](https://github.com/git-town/git-town/pull/2703))
  - `perennial-branch-names` is now `perennial-branches` ([#2702](https://github.com/git-town/git-town/pull/2702))
  - `sync-strategy` is now `sync-feature-strategy` ([#2697](https://github.com/git-town/git-town/pull/2697))
  - `pull-branch-strategy` is now `sync-perennial-strategy` ([#2693](https://github.com/git-town/git-town/pull/2693))
- `git ship` by default no longer syncs the branch to ship. Set the `sync-before-ship` flag to restore the old behavior. This allows shipping only when the tests pass ([#2735](https://github.com/git-town/git-town/pull/2735)).
- Creating proposals on BitBucket uses an updated URL ([#2692](https://github.com/git-town/git-town/pull/2692)).
- `git town config reset` now also deletes the branch lineage. This helps get you out of more configuration snafus ([#2733](https://github.com/git-town/git-town/pull/2733)).

#### New Features

- The new `sync-before-ship` config option prevents `git ship` from updating the branch it is about to ship. The old behavior makes sense when shipping branches locally but is conflicting with the requirements for tests to pass on CI before shipping via the hosting platform ([#2714](https://github.com/git-town/git-town/pull/2714)).

#### Bug Fixes

- allow renaming local-only branches ([#2710](https://github.com/git-town/git-town/pull/2710))
- `git repo` and `git propose` always open the browser page at the default port, even if the `origin` remote points to a custom port ([#2730](https://github.com/git-town/git-town/pull/2730))

## 10.0.3 (2023-11-25)

#### Bug Fixes

- Fix killing perennial branches ([#2679](https://github.com/git-town/git-town/pull/2679))

## 10.0.2 (2023-11-7)

#### Bug Fixes

- Fix wrong error message when `status.short` is enabled in Git config ([#2650](https://github.com/git-town/git-town/pull/2650))

## 10.0.1 (2023-11-2)

#### Bug Fixes

- Fix crash if commits contain "[" ([#2645](https://github.com/git-town/git-town/pull/2645))

## 10.0.0 (2023-10-27)

Git Town 10 improves support for shipping branches via the code hosting web UI instead of running `git ship`. After merging your branches remotely, run `git sync --all` to sync all local branches and remove the ones shipped at the remote. Don't worry, Git Town ensures that branches which contain unshipped changes won't get deleted. `git undo` brings deleted branches back.

Git Town 10 has improved performance, robustness, and reliability thanks to a large-scale modernization of the Git Town's architecture. Git Town now runs fewer Git commands under the hood to investigate the state of your Git repository. `git undo` now works for all commands thanks to a new undo engine that diffs the before and after state of your Git repo.

Git Town 10 starts a larger effort to remove redundant commands and make Git Town's configuration options more consistent and intuitively named.

#### BREAKING CHANGES

- `git sync` now also removes local branches with a deleted tracking branch, after verifying that those local branches contain no unshipped changes ([#2038](https://github.com/git-town/git-town/pull/2038))
- `git town prune-branches` has been sunset, run `git sync` instead ([#2579](https://github.com/git-town/git-town/pull/2579))
- Git Town's statefile on disk has a new format, you might have to run `git town status reset` to avoid runtime errors ([#2446](https://github.com/git-town/git-town/pull/2446))
- `git ship` no longer ships branches that exist solely at the remote. Moving forward branches to ship must exist on your local machine. Use the web UI of your code hosting service to ship remote branches. ([#2367](https://github.com/git-town/git-town/pull/2367), [#2372](https://github.com/git-town/git-town/pull/2372))
- `git kill` no longer deletes branches that exist solely at the remote. Delete them by running `git push origin :branchname` or via the web UI of your code hosting service ([#2368](https://github.com/git-town/git-town/pull/2368))
- `git hack` no longer has the `-p` option. Use `git append` and `git prepend` instead ([#2577](https://github.com/git-town/git-town/pull/2577))
- Git Town no longer considers it an error if there is nothing to abort or continue ([#2631](https://github.com/git-town/git-town/pull/2631), [#2632](https://github.com/git-town/git-town/pull/2632))
- querying the version of the installed Git Town binary is now compatible with the way Git does it: `git-town --version` instead of `git-town version` ([#2603](https://github.com/git-town/git-town/pull/2603))
- v10 renames the `debug` parameter to `verbose` because all it does is print more information ([#2598](https://github.com/git-town/git-town/pull/2598))
- updated GitLab support, please report regressions

#### New Features

- support for running Git Town on computers that use different language than English ([#2478](https://github.com/git-town/git-town/pull/2478))
- `git undo` works for all commands now ([#2484](https://github.com/git-town/git-town/pull/2484))
- CLI output now contains requests to the code hosting API ([#2340](https://github.com/git-town/git-town/pull/2340))
- CLI output now describes changes the branch ancestry ([#2558](https://github.com/git-town/git-town/pull/2558))
- `git town switch` now displays the output of the command to switch branches ([#2602](https://github.com/git-town/git-town/pull/2602))
- environment variables now override all GitHub API operations ([#2593](https://github.com/git-town/git-town/pull/2593))
- community-contributed installation for BSD via FreshPorts ([#2553](https://github.com/git-town/git-town/pull/2553))
- less force-deleting of branches ([#2539](https://github.com/git-town/git-town/pull/2539))

#### Bug Fixes

- fix broken version number in release binaries ([#2333](https://github.com/git-town/git-town/pull/2333))
- fix crash when a configured branch parent is empty ([#2626](https://github.com/git-town/git-town/pull/2626))
- fix crash when running `set-parent` on large monorepos ([#2623](https://github.com/git-town/git-town/pull/2623))
- when deleting perennial branches, remove the ancestry information of their children ([#2540](https://github.com/git-town/git-town/pull/2540))

## 9.0.1 (2023-07-29)

Git Town should now run a bit faster because it runs fewer Git commands under the hood.

#### Bug Fixes

- Fix for missing `UpdateProposalTargetStep` ([#2288](https://github.com/git-town/git-town/pull/2288))
- Print statistics when removing aliases ([#2325](https://github.com/git-town/git-town/pull/2325))
- Fix broken version information in release binaries ([#2333](https://github.com/git-town/git-town/pull/2333))

## 9.0.0 (2023-04-07)

#### BREAKING CHANGES

Git Town 9.0 supports the new API URLs that become official in GitLab v16. If you use an older version of GitLab, Git Town's integration with GitLab's API might no longer work. The fix is to update your GitLab installation to at least v15 ([#2249](https://github.com/git-town/git-town/pull/2249))

#### New Features

- read the token to use for the GitHub API from environment variables `GITHUB_TOKEN` or `GITHUB_AUTH_TOKEN` in addition to the already existing option to store it in the Git configuration ([#2217](https://github.com/git-town/git-town/pull/2217))

## 8.0.0 (2023-04-07)

Some ergonomics improvements that change existing command names, hence the major version bump. If you use the shorter aliases for Git Town commands, please run `git town aliases add` after updating to v8.0.

#### BREAKING CHANGES

- rename `new-branch-push-flag` command to `push-new-branches` ([#1980](https://github.com/git-town/git-town/pull/1980))
- all commands that display/update configuration are now subcommands of the `config` command ([#1963](https://github.com/git-town/git-town/pull/1963), [#1976](https://github.com/git-town/git-town/pull/1976))
- all commands that help install Git Town are now subcommands of the `install` command ([#1969](https://github.com/git-town/git-town/pull/1969))
- moves the default `git town completion` and `git town completions` commands under `git town install completions` ([#1969](https://github.com/git-town/git-town/pull/1969), [#1970](https://github.com/git-town/git-town/pull/1970))
- installation of the shorter command aliases changes from `git town alias true` to `git town install aliases add` and `... remove` ([#1965](https://github.com/git-town/git-town/pull/1965), [#1966](https://github.com/git-town/git-town/pull/1966), [#1968](https://github.com/git-town/git-town/pull/1968))
- renames the `push-verify` configuration option to `push-hook` ([#1989](https://github.com/git-town/git-town/pull/1989))
- automatically renames old `push-verify` configuration settings to the new `push-hook` ([#2209](https://github.com/git-town/git-town/pull/2209))
- rename `git set-parent-branch` to `git set-parent` ([#2114](https://github.com/git-town/git-town/pull/2114))
- stores the runstate in the platform-specific config directory (`~/.config`) instead of the global temp dir ([#2126](https://github.com/git-town/git-town/pull/2126))

#### New Features

- the new `git town switch` command allows switching branches via a UI that visualizes the branch hierarchy ([#2106](https://github.com/git-town/git-town/pull/2106), [#2108](https://github.com/git-town/git-town/pull/2108))
- aliases `git town diff-parent` to `git diff-parent` ([#2128](https://github.com/git-town/git-town/pull/2128))
- accepts more formats for boolean configuration values like "true", "yes", "on", "t", "1" ([#1978](https://github.com/git-town/git-town/pull/1978), [#1979](https://github.com/git-town/git-town/pull/1979))
- configuration command to set/display the `push-hook` config setting ([#1991](https://github.com/git-town/git-town/pull/1991))

## 7.9.0 (2023-01-22)

#### New Features

- rebase feature branches against their parent branch using the new [sync-strategy option](https://www.git-town.com/preferences/sync-strategy.html) ([#1950](https://github.com/git-town/git-town/pull/1950))
  - configure using `git town sync-strategy (merge | rebase)`
- disable Git's `pre-push` hook using the new `push-verify` option ([#1958](https://github.com/git-town/git-town/pull/1958))
- support for [GitLab subgroups](https://docs.gitlab.com/ee/user/group/subgroups) ([#1943](https://github.com/git-town/git-town/pull/1943))

#### Bug Fixes

- support GitLab SaaS repos whose name contains "gitlab" ([#1926](https://github.com/git-town/git-town/pull/1926))

## 7.8.0 (2022-08-07)

#### New Features

- update Regex for hostname extraction to support more ssh usernames ([#1883](https://github.com/git-town/git-town/pull/1883))
- merge GitLab merge requests when shipping ([#1874](https://github.com/git-town/git-town/pull/1874))

#### Bug Fixes

- fix tests on non-English locales ([#1875](https://github.com/git-town/git-town/pull/1875))
- fix bug in undo of "git hack" in local repo ([#1804](https://github.com/git-town/git-town/pull/1804))

## 7.7.0 (2022-01-22)

#### New Features

- support for Apple Silicon ([#1735](https://github.com/git-town/git-town/pull/1735), [#1736](https://github.com/git-town/git-town/pull/1736))
- ignore changes in submodules during sync ([#1744](https://github.com/git-town/git-town/pull/1744))
- improved CLI interface including better shell autocompletions ([#1722](https://github.com/git-town/git-town/pull/1722))
- shell-based installer for *nix systems ([#1707](https://github.com/git-town/git-town/pull/1707))
- new website ([#1684](https://github.com/git-town/git-town/pull/1684))
- Make command shows dependency tree within the codebase ([#1725](https://github.com/git-town/git-town/pull/1725))
- Go API now has the same major version number as the binary ([#1677](https://github.com/git-town/git-town/pull/1677))

#### Bug Fixes

- fix author not set properly ([1686](https://github.com/git-town/git-town/pull/1686))
- filenames of assets at GitHub releases are all lowercase now ([#1710](https://github.com/git-town/git-town/pull/1710))
- option for more compatible shell autocompletion without descriptions ([#1493](https://github.com/git-town/git-town/pull/1493))

## 7.6.0 (2021-11-23)

#### New Features

- print diagnostic information on command failure ([#1667](https://github.com/git-town/git-town/pull/1667))

## 7.5.0 (2021-03-10)

#### New Features

- add log after command that causes auto abort ([#1635](https://github.com/git-town/git-town/pull/1635))

#### Bug Fixes

- fix panic when continuing rebase ([#1615](https://github.com/git-town/git-town/pull/1615))
- fix panic when continuing a command that includes the fetch upstream step ([#1617](https://github.com/git-town/git-town/pull/1617))
- fix GitHub API log when shipping with the GitHub driver ([#1622](https://github.com/git-town/git-town/pull/1622))
- fix panic when aborting a command that includes discard open changes step ([#1631](https://github.com/git-town/git-town/pull/1631))
- fix continuing sync from subfolder ([#1637](https://github.com/git-town/git-town/pull/1637))

## 7.4.0 (2020-07-05)

Version 7.4.0 sports a vastly overhauled internal architecture that provides more robust error handling, improved error messages, and a much better developer experience. Ruby is no longer a development dependency.

#### New Features

- improved installation experience: MSI installer for Windows, `.deb` and `.rpm` packages for Linux, archives with properly named binaries for all other use cases ([#1589](https://github.com/git-town/git-town/pull/1589))
- "diff-parent" command ([#1385](https://github.com/git-town/git-town/pull/1518))
- support for Gitea hosting service ([#1518](https://github.com/git-town/git-town/pull/1518))
- print URLs to open when browser is not available ([#1318](https://github.com/git-town/git-town/pull/1318))
- autocompletion for bash, zsh, fish, powershell ([#1492](https://github.com/git-town/git-town/pull/1492))
- list parent configurations for branches that are children of a branch that does not have its parent configured ([#1436](https://github.com/git-town/git-town/pull/1436))

#### Bug Fixes

- improved error messages

## 7.3.0 (2019-11-05)

#### New Features

- add option to disable auto sync upstream

#### Bug Fixes

- update docs for code-hosting config

## 7.2.1 (2019-05-06)

#### Bug Fixes

- prune branches now properly updates perennial branch config
- support branch names with special characters
- fix the prompt on Windows CMD terminals
- clear the runstate after undo to prevent running `git town undo` twice
- fix Fish shell autocomplete
- fix hosting service naming

## 7.2.0 (2018-06-01)

#### New Features

- `git town config`: print perennial branch trees
- `git town hack`: add `-p` option which prompts for the parent branch instead of using the main branch
- when fetching the origin repository, fetch tags that are not attached to pulled commits

#### Changes

- fetch only the main branch when fetching the upstream repository

## 7.1.1 (2018-04-09)

#### Bug Fixes

- strip colors from the output of git commands run internally. This caused errors if you had git configured with `color.ui=always`

## 7.1.0 (2018-04-05)

#### New Features

- automatically remove outdated configuration

## 7.0.0 (2018-04-03)

#### BREAKING CHANGES

- `git town config`: `reset` and `setup` are now subcommands instead of flags
- `--abort`, `--continue`, `--skip`, `--undo` flags removed. Instead there are now top level commands `git town abort`, `git town continue`, `git town skip`, `git town undo`

#### New Features

- Catches when there is an unfinished state from a git town command that hit conflicts. If you try to run another git town command, Git Town will prompt you on how to resolve the unfinished state. The unfinished state can be discarded and there is also a new top level command `git town discard` to delete the state of the last run command.

#### Bug Fixes

- skip perennial branch prompt if there are no options

## 6.0.2 (2018-01-26)

#### Bug Fixes

- fix parsing of git config when a value contains a newline

## 6.0.1 (2018-01-24)

#### Bug Fixes

- fixes displayed version number

## 6.0.0 (2018-01-15)

#### BREAKING CHANGES

- `git town set-parent-branch`: update interface to no longer accept arguments and instead prompt the user for the parent of the current branch
- `git town perennial-branches`: update the interface to add / remove perennial branches. Run `git town perennial-branch update` to receive the same prompt as initial configuration.
- Rename `hack-push-flag` to `new-branch-push-flag`. Please reconfigure if you are not using the default.

#### New Features

- `git town new-branch-push-flag`: add `--global` flag to set your default value. Any locally configured value will override.
- add `--verbose` flag to see all the git commands runs under the hood
- speed improvement from reducing the number of git commands run under the hood

## 5.1.0 (2017-12-05)

#### New Features

- Nicer prompts from https://github.com/AlecAivazis/survey
- Parent branch prompt: add option to make the branch a perennial branch

#### Bug Fixes

- `git ship`: fix bug when encountering a merge conflict and using a code hosting driver ([#1060](https://github.com/git-town/git-town/issues/1060))

## 5.0.0 (2017-08-16)

#### BREAKING CHANGES

- `git new-pull-request / repo`: support for ssh identities changed
  - Before: ssh identity needed to include "github", "gitlab" or "bitbucket"
  - Now: Run `git config git-town.code-hosting-origin-hostname <hostname>` where hostname matches what is in your ssh config file

#### New Features

- `git new-pull-request / repo`: support for self hosted versions
  - Run `git config git-town.code-hosting-driver <driver>` where driver is "bitbucket", "github", or "gitlab"
- `git sync`: add `--dry-run` flag to view the planned commands without running them
- `git ship`: when merging via the GitHub API, update the default commit message to include the PR title and number

## 4.2.1 (2017-08-16)

#### Bug Fixes

- add missing dependency to vendor folder (required for building on Homebrew)

## 4.2.0 (2017-08-15)

#### New Features

- Update all commands to support offline mode (lack of an internet connection)
  - Display / update offline mode with `git town offline [(true | false)]`
- `git ship`
  - add ability to ship hotfixes to perennial branches
  - add ability to merge via GitHub API when applicable. See [documentation](website/src/commands/ship.md) for more info.

## 4.1.2 (2017-06-08)

#### Bug Fixes

- temporary file: use operating system temporary directory instead of hardcoding `/tmp`

## 4.1.1 (2017-06-07)

#### Bug Fixes

- temporary file: make parent directories if needed ([#955 comment](https://github.com/git-town/git-town/issues/955#issuecomment-306041043))

## 4.1.0 (2017-06-01)

#### New Features

- `git new-pull-request`, `git repo`: support more commands to open browsers (`cygstart`, `x-www-browser`, `firefox`, `opera`, `mozilla`, `netscape`)
- Add longer descriptions for commands which appear when running `git town help <command>` or `git town <command> --help`

#### Changes

- make `hack-push-flag` false by default (was true before) ([#929](https://github.com/git-town/git-town/issues/929))

#### Bug Fixes

- replace all non-alpha numeric characters in temporary filename ([#925](https://github.com/git-town/git-town/issues/925))
- fix spacing in parent branch prompts
- enforce Git version 2.7.0 or higher

## 4.0.1 (2017-05-21)

#### Bug Fixes

- fix infinite loop when prompting for parent branch and there are perennial branches configured
- enforce a Git version 2.6.0 or higher
- fix `ship` when the supplied branch is equal to the current branch and there are open changes
- allow running `alias` command in non-git directories

## 4.0.0 (2017-05-12)

#### BREAKING CHANGES

- rewrite in go, Git Town is now a single, stand-alone binary
  - first-class Windows support
  - This breaks existing aliases. If you have the default aliases setup, reconfigure them with `git town alias true`

## 3.1.0 (2017-03-27)

#### New Features

- `git new-pull-request`, `git repo`:
  - support `ssh://` urls (thanks to @zhangwei)
  - add GitLab support (thanks to @dgjnpr)

## 3.0.0 (2017-02-07)

#### BREAKING CHANGES

- `git hack`: no longer accepts a parent branch (functionality moved to `git append`)

#### New Features

- `git append`: create a new branch as a child of the current branch
- `git prepend`: create a new branch as a parent of the current branch
- `git rename-branch`: implicitly uses the current branch if only one branch name provided

#### Bug Fixes

- fix incorrectly reported branch loop ([#785](https://github.com/git-town/git-town/issues/785))

## 2.1.0 (2016-12-26)

#### New Features

- support SSH identities ([#739](https://github.com/git-town/git-town/issues/739))

#### Bug Fixes

- update stashing strategy to avoid use of `git stash -u` which can delete ignored files ([#744](https://github.com/git-town/git-town/issues/744))
- fix merge conflicts resolution that results in no changes ([#753](https://github.com/git-town/git-town/issues/753))
- `git hack`: prompt for parent branch if unknown ([#760](https://github.com/git-town/git-town/issues/760))
- prevent parent branch loops ([#751](https://github.com/git-town/git-town/issues/751))

## 2.0.0 (2016-09-18)

#### BREAKING CHANGES

- All commands now have a `town-` prefix. Example `git town-sync`. This is to prevent conflicts with `git-extras` which adds git commands by the same name and `hub` which wants you to alias git to it and adds commands by the same name.
  - Use [git aliases](https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases) to remove the `town-` prefix if you would like. Run `git town alias true` to add aliases for all `git-town` commands (skips commands which would overwrite existing aliases).

## 1.0.0 (2016-08-05)

#### New Features

- `git town set-parent-branch <child_branch_name> <parent_branch_name>`: to update a parent branch ([#729](https://github.com/git-town/git-town/issues/729))

#### Bug Fixes

- `git sync --all`: don't prompt for parent of perennial branches ([#727](https://github.com/git-town/git-town/issues/727))

## 0.10.1 (2016-06-23)

#### New Features

- `git hack`: add configuration to omit pushing new branches ([#720](https://github.com/git-town/git-town/issues/720))

#### Bug Fixes

- configuration: make branch order consistent
- `git ship`: update uncommitted changes error message ([#718](https://github.com/git-town/git-town/issues/718))

## 0.10.0 (2016-01-21)

#### BREAKING CHANGES

- `git prune-branches`: new functionality - delete branches whose tracking branch no longer exists ([#702](https://github.com/git-town/git-town/issues/702))

#### New Features

- `git sync`: add configuration option to merge the main branch / perennial branches with their upstream ([#671](https://github.com/git-town/git-town/issues/671))
- `git hack`, `git ship`: support for running in subfolders

#### Bug Fixes

- internationalize check for undefined function ([#678](https://github.com/git-town/git-town/issues/678))
- `git new-pull-request`: ability to continue after conflicts ([#700](https://github.com/git-town/git-town/issues/700))

## 0.9.0 (2015-10-17)

#### BREAKING CHANGES

- remove `git sync-fork`

#### New Features

- `git new-pull-request`: support forked repos ([#682](https://github.com/git-town/git-town/issues/682))
- `git sync`: if there is a remote upstream, syncs the main branch with its upstream counterpart ([#685](https://github.com/git-town/git-town/issues/685))

## 0.8.0 (2015-10-14)

#### BREAKING CHANGES

- removed `git extract`
- update internal storage of perennial branches
  - if you have configured more than one perennial branch, you will need to reset your configuration
    - `git town config --reset`
    - `git town config --setup` or follow the prompt the next time you run a Git Town command

#### New Features

- configuration prompt: allow user to select branch by number, ability to recover from bad input
- parent branch prompt: show description and branch list once per command
- preserve checkout history so that `git checkout -` works as expected alongside Git Town commands ([#65](https://github.com/git-town/git-town/issues/65))
- `git hack`: pushes the new branch to the remote repository ([#664](https://github.com/git-town/git-town/issues/664))
- `git new-pull-request`: syncs the branch before creating the pull request ([#367](https://github.com/git-town/git-town/issues/367))
- `git sync --all`: pushes tags ([#464](https://github.com/git-town/git-town/issues/464))
- `git town config`: shows branch ancestry ([#651](https://github.com/git-town/git-town/issues/651))

#### Bug Fixes

- `git town version`: Homebrew installs no longer print date and SHA ([#631](https://github.com/git-town/git-town/issues/631))

## 0.7.3 (2015-09-02)

- `git kill`: can remove branches that exist on the remote and not locally ([#380](https://github.com/git-town/git-town/issues/380))
- `git ship`: prompt when there is more than one author ([#486](https://github.com/git-town/git-town/issues/486))

## 0.7.2 (2015-08-28)

- `git sync --all`: fix parent branch prompt
- `git ship`: comment out default commit message ([#382](https://github.com/git-town/git-town/issues/382))

## 0.7.1 (2015-08-27)

- `git ship`: add ability to ship remote branches ([#593](https://github.com/git-town/git-town/issues/593))
- `git new-pull-request`/`git repo`: remove empty line output ([#602](https://github.com/git-town/git-town/issues/602))
- `git kill`: prompt for unknown parent branch ([#603](https://github.com/git-town/git-town/issues/603))
- `git sync --all`: prompt for unknown parent branch ([#604](https://github.com/git-town/git-town/issues/604))
- support branch names with forward slashes (along with any valid branch name) ([#608](https://github.com/git-town/git-town/issues/608))

## 0.7.0 (2015-08-24)

- fix `git ship --undo` ([#550](https://github.com/git-town/git-town/issues/550))
- rename `non-feature-branches` to `perennial-branches` ([#344](https://github.com/git-town/git-town/issues/344))
  - configuration is automatically updated to support this
- support for nested feature branches ([#529](https://github.com/git-town/git-town/issues/529))
- add `git rename-branch` ([#474](https://github.com/git-town/git-town/issues/474))
- rename `git pull-request` to `git new-pull-request` ([#413](https://github.com/git-town/git-town/issues/413), [#507](https://github.com/git-town/git-town/issues/507))
- add SHA and date to output of `git town version` for manual installs
- show error when trying to continue after a successful command ([#364](https://github.com/git-town/git-town/issues/364))

## 0.6.0 (2015-04-02)

- support for working without a remote repository for **git extract**, **git hack**, **git kill**, **git ship**, and **git sync**
  - implemented by our newest core committer @ricmatsui
- **git pr** renamed to **git pull-request**
  - set up an alias with `git config --global alias.pr pull-request`
- **git ship**
  - now accepts all `git commit` options
  - author with the most commits is automatically set as the author (when not the committer) ([#335](https://github.com/git-town/git-town/issues/335))
- **git pr/repo**
  - improved linux compatibility by trying `xdg-open` before `open`
- improved error messages when run outside a git repository
- improved setup wizard for initial configuration in a git repository
- added [contribution guide](CONTRIBUTING.md)
- added tutorial

## 0.5.0 (2015-01-08)

- Manual installs need to update their `PATH` to point to the `src` folder within their clone of the repository
- **git extract:**
  - errors if branch exists remotely ([#236](https://github.com/git-town/git-town/issues/236))
  - removed restriction: need to be on a feature branch ([#269](https://github.com/git-town/git-town/issues/269))
  - added restriction: errors if the current branch does not have any have extractable commits and the user provided no commits (commits not in the main branch) ([#269](https://github.com/git-town/git-town/issues/269))
- **git hack:** errors if branch exists remotely ([#237](https://github.com/git-town/git-town/issues/237))
- **git kill:**
  - optional branch name ([#126](https://github.com/git-town/git-town/issues/126))
  - does not error if tracking branch is already deleted ([#196](https://github.com/git-town/git-town/issues/196))
- **git pr:**
  - linux compatibility ([#232](https://github.com/git-town/git-town/issues/232))
  - compatible with more variants of specifying a Bitbucket or GitHub remote ([#271](https://github.com/git-town/git-town/issues/271))
  - compatible with respository names that contain ".git" ([#305](https://github.com/git-town/git-town/issues/305))
- **git repo:** view the repository homepage ([#140](https://github.com/git-town/git-town/issues/140))
- **git sync:**
  - `--all` option to sync all local branches ([#83](https://github.com/git-town/git-town/issues/83))
  - abort with correct state after main branch updates and tracking branch conflicts ([#228](https://github.com/git-town/git-town/issues/228))
- **git town**: view and change Git Town configuration and improved view help page ([#98](https://github.com/git-town/git-town/issues/98))
- auto-completion for [Fish shell](https://fishshell.com) ([#177](https://github.com/git-town/git-town/issues/177))

## 0.4.1 (2014-12-02)

- **git pr:** create a new pull request ([#138](https://github.com/git-town/git-town/issues/138), [40d22e](https://github.com/git-town/git-town/commit/40d22eb1703ac96a58ac5052e70d20d7bdb9ac73))
- **git ship:**
  - empty commit message aborts the command ([#153](https://github.com/git-town/git-town/issues/153), [0bc84e](https://github.com/git-town/git-town/commit/0bc84ee626299896661fe1754cfa227630725bb9))
  - abort when there are no shippable changes ([#188](https://github.com/git-town/git-town/issues/188), [52fd94](https://github.com/git-town/git-town/commit/52fd94eca05bd3c2db5e7ac36121f08e56b9558b))
- **git sync:**
  - can now continue after resolving conflicts (no need to commit or continue rebasing) ([#123](https://github.com/git-town/git-town/issues/123), [1a50ad](https://github.com/git-town/git-town/commit/1a50ad689a752f4eaed663e0ab22184621ee96a2))
  - restores deleted tracking branch ([#165](https://github.com/git-town/git-town/issues/165), [259464](https://github.com/git-town/git-town/commit/2594646ad853d83a6d697354d66755a374e42b8a))
- **git extract:** errors if branch already exists ([#128](https://github.com/git-town/git-town/issues/128), [75f498](https://github.com/git-town/git-town/commit/75f498771f19326f03bd1fd1bb70c9d9851b53f3))
- **git sync-fork:** no longer automatically sets upstream configuration ([865030](https://github.com/git-town/git-town/commit/8650301a3ea40a989562a991960fa0d41b26f7f7))
- remove needless checkouts for **git-ship**, **git-extract**, and **git-hack** ([#150](https://github.com/git-town/git-town/issues/150), [#155](https://github.com/git-town/git-town/issues/155), [8b385a](https://github.com/git-town/git-town/commit/8b385a745cf7ed28638e0a5c9c24440b7010354c), [35de43](https://github.com/git-town/git-town/commit/35de43156d9c6092840cd73456844b90acc36d8e))
- linters for shell scripts and ruby tests ([#149](https://github.com/git-town/git-town/issues/149), [076668](https://github.com/git-town/git-town/commit/07666825b5d60e15de274746fc3c26f72bd7aee2), [651c04](https://github.com/git-town/git-town/commit/651c0448309a376eee7d35659d8b06f709b113b5))
- rake tasks for development ([#170](https://github.com/git-town/git-town/issues/170), [ba74cf](https://github.com/git-town/git-town/commit/ba74cf30c8001941769dcd70410dbd18331f2fe9))

## 0.4.0 (2014-11-13)

- **git kill:** removes a feature branch ([#87](https://github.com/git-town/git-town/issues/87), [edd7d8](https://github.com/git-town/git-town/commit/edd7d8180eb76717fd72e77d2c75edf8e3b6b6ca))
- **git sync:** pushes tags to the remote when running on the main branch ([#68](https://github.com/git-town/git-town/issues/68), [71b607](https://github.com/git-town/git-town/commit/71b607988c00e6dfc8f2598e9b964cc2ed4cfc39))
- **non-feature branches:** are not shipped and do not merge main when syncing ([#45](https://github.com/git-town/git-town/issues/45), [31dce1](https://github.com/git-town/git-town/commit/31dce1dfaf11e1e17f17e141a26cb38360ab731a))
- **git ship:**
  - merges main into the feature branch before squash merging ([#61](https://github.com/git-town/git-town/issues/61), [82d4d3](https://github.com/git-town/git-town/commit/82d4d3e745732cb397850a4c047826ba485e2bdb))
  - errors if the feature branch is not ahead of main ([#86](https://github.com/git-town/git-town/issues/86), [a0ace5](https://github.com/git-town/git-town/commit/a0ace5bb5e992c193df8abe4b0aca984c302c323))
  - git ship takes an optional branch name ([#95](https://github.com/git-town/git-town/issues/95), [cbf020](https://github.com/git-town/git-town/commit/cbf020fc3dd6d0ce49f8814a92f103e243f9cd2b))
- updated output to show each git command and its output, updated error messages ([8d8973](https://github.com/git-town/git-town/commit/8d8973aaa58394a123ceed2811271606f4e1aaa9), [60e1d8](https://github.com/git-town/git-town/commit/60e1d8299ebbb0e75bdae057e864d17e1f9a3ce7), [408e69](https://github.com/git-town/git-town/commit/408e699e5bdd3af524b2ea64669b81fea3bbe60b))
- skips unnecessary pushes ([0da896](https://github.com/git-town/git-town/commit/0da8968aef29f9ecb7326e0fafb5976f51789dca))
- **man pages** ([609e11](https://github.com/git-town/git-town/commit/609e11400818604328885df86c02ee4630410e12), [164f06](https://github.com/git-town/git-town/commit/164f06bc8bf00d9e99ce0416f408cf62959dc833), [27b257](https://github.com/git-town/git-town/commit/27b2573ca5ffa9ae7930f8b5999bbfdd72bd16d9))
- **git prune-branches** ([#48](https://github.com/git-town/git-town/issues/48), [7a922e](https://github.com/git-town/git-town/commit/7a922ecd9e03d20ed5a0c159022e601cebc80313))
- **Cucumber:** optional Fuubar output ([7c5402](https://github.com/git-town/git-town/commit/7c540284cf46bd49a7623566c1343285813524c6))

## 0.3 (2014-10-10)

- multi-user support for feature branches ([#35](https://github.com/git-town/git-town/issues/35), [ca0882](https://github.com/git-town/git-town/commit/ca08820c68457bddf6b8fff6c3ef3d430b905d9b))
- **git sync-fork** ([#22](https://github.com/git-town/git-town/issues/22), [1f1f9f](https://github.com/git-town/git-town/commit/1f1f9f98ffa7288d6a5982ec0c9e571695590fe1))
- stores configuration in the Git configuration instead of a dedicated file ([8b8695](https://github.com/git-town/git-town/commit/8b86953d7c7c719f28dbc7af6e86d02adaf2053e))
- removes redundant fetches from the central repo per session ([#15](https://github.com/git-town/git-town/issues/15), [43400a](https://github.com/git-town/git-town/commit/43400a5b968a47eb55484f73e34026f66b1e939a))
- automatically prunes remote branches when fetching updates ([86100f](https://github.com/git-town/git-town/commit/86100f08866f19a0f4e80f470fe8dcc6996ddc2c))
- always cleans up abort and continue scripts after using one of them ([3be4c0](https://github.com/git-town/git-town/commit/3be4c06635a943f378287963ba30e4306fcd9802))
- simpler readme, dedicated RDD document
- **<a href="https://cukes.info" target="_blank">Cucumber</a>** feature specs (you need Ruby 2.x) ([c9d175](https://github.com/git-town/git-town/commit/c9d175fe2f28fbda3f662454f54ed80306ce2f46))
- much faster testing thanks to fully local test Git repos ([#25](https://github.com/git-town/git-town/issues/25), [c9d175](https://github.com/git-town/git-town/commit/c9d175fe2f28fbda3f662454f54ed80306ce2f46))

## 0.2.2 (2014-06-10)

- fixes "unary" error messages
- lots of output and documentation improvements

## 0.2.1 (2014-05-31)

- better terminal output
- Travis CI improvements
- better documentation

## 0.2.0 (2014-05-29)

- displays the duration of specs
- pulls the main branch only if it has a remote
- --abort options to abort failed Git Town operations
- --continue options to continue some Git Town operations after fixing the underlying issues
- installation through Homebrew
- colored test output
- display summary after tests
- exit with proper status codes
- better documentation

## 0.1.0 (2014-05-22)

- git hack, git sync, git extract, git ship
- basic test framework
- Travis CI integration
- self-hosting: uses Git Town for Git Town development
