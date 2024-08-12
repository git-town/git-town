# Git Town Changelog

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
- added [contribution guide](/docs/CONTRIBUTING.md)
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
