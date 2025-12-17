# git town sync

<a type="git-town-command" />

```command-summary
git town sync [-a | --all] [--(no)-auto-resolve] [-d | --(no)-detached] [--dry-run] [--gone] [-h | --help] [-p | --prune] [--(no)-push] [-s | --stack] [-v | --verbose]
```

The _sync_ command ("synchronize this branch") updates your local Git workspace
with what happened in the rest of the repository.

You can (and should) sync all branches many times per day without thinking about
it, even in the middle of ongoing work. If a sync goes wrong, you can safely go
back to the exact state you repo was in before the sync by running
[git town undo](undo.md).

- pulls and pushes updates from all parent branches and the tracking branch
- deletes branches whose tracking branch was deleted at the remote if they
  contain no unshipped changes
- removes commits of deleted branches from their descendent branches, unless
  when using the
  [merge sync strategy](../preferences/sync-feature-strategy.md#merge).
- safely stashes away uncommitted changes and restores them when done
- does not pull, push, or merge depending on the configured
  [branch type](../branch-types.md)

If the parent branch is not known, Git Town looks for a pull/merge request for
this branch and uses its parent branch. Otherwise it prompts you for the parent.

### Sync frequently

Merge conflicts are not fun and can break code. Minimize them by making it a
habit to sync your branches regularly and frequently. When properly configured,
`git town sync --all` can synchronize all your local branches the right way
without losing changes, even in edge cases.

If you don't sync because:

- you don't want to pull in new changes from the main branch:
  [sync detached](sync.md#-d--detached--no-detached).
- you don't want to increase pressure on your CI server:
  [sync without pushing](sync.md#--push--no-push) or indicate in your commit
  messages to CI to skip test runs
  - [BitBucket](https://support.atlassian.com/bitbucket-cloud/kb/how-to-skip-triggering-an-automatic-pipeline-build-using-skip-ci-label)
  - [Gitea](https://docs.gitea.com/administration/config-cheat-sheet#actions-actions)
  - [GitHub](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/skip-workflow-runs)
  - [GitLab](https://docs.gitlab.com/ci/pipelines/#skip-a-pipeline)
  - [Forgejo](https://forgejo.org/docs/latest/admin/config-cheat-sheet/#actions-actions)

### Why does Git Town sometimes not sync the tracking or parent branch?

Git Town detects whether there are any changes that need to be synced, and might
skip unnecessary sync operations that wouldn't produce any changes.

### Why does Git Town sometimes update a local branch whose tracking branch was deleted before deleting it?

If a remote branch was deleted at the remote, it is considered obsolete and
`git town sync` will remove its local counterpart. To guarantee that this
doesn't lose unshipped changes in the local branch, `git town sync` needs to
prove that the branch to be deleted contains no unshipped changes.

The easiest way to prove that is when the local branch was in sync with its
tracking branch before Git Town runs `git fetch`. This is another reason to run
`git town sync` regularly.

If a local shipped branch is not in sync with its tracking branch on your
machine, Git Town must check for unshipped local changes by diffing the branch
to delete against its parent branch. Only branches with an empty diff can be
deleted safely. For this to work, Git Town needs to sync the branch first, even
if it's going to be deleted right afterwards.

## Options

#### `-a`<br>`--all`

By default this command syncs only the current branch. The `--all` aka `-a`
parameter makes Git Town sync all local branches.

#### `--auto-resolve`<br>`--no-auto-resolve`

Disables automatic resolution of
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts).

#### `-d`<br>`--detached`<br>`--no-detached`

The `--detached` aka `-d` flag enables
[detached mode](../preferences/detached.md) for the current command. If detached
mode is enabled through [configuration data](../preferences/detached.md), the
`--no-detached` flag disables detached mode for the current command.

In detached mode, feature branches don't receive updates from the perennial
branch at the root of your branch hierarchy. This can be useful in busy
monorepos.

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `--gone`

Sync all local branches whose remote is gone. This effectively removes all local
branches that were shipped or deleted at the remote.

#### `-h`<br>`--help`

Display help for this command.

#### `-p`<br>`--prune`

The `--prune` aka `-p` flag removes (prunes) empty branches, i.e. branches that
effectively don't make any changes.

#### `--push`<br>`--no-push`

The `--push`/`--no-push` argument overrides the
[push-branches](../preferences/push-branches.md) config setting.

#### `-s`<br>`--stack`

The `--stack` aka `-s` parameter makes Git Town sync all branches in the stack
that the current branch belongs to.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## Configuration

[sync-perennial-strategy](../preferences/sync-perennial-strategy.md) configures
whether perennial branches merge their tracking branch or rebase against it.

[sync-feature-strategy](../preferences/sync-feature-strategy.md) configures
whether feature branches merge their parent and tracking branches or rebase
against them.

If the repository contains a Git remote called `upstream` and the
[sync-upstream](../preferences/sync-upstream.md) setting is enabled, Git Town
also pulls new commits from the upstream's main branch.

[sync-tags](../preferences/sync-tags.md) configures whether Git Town syncs Git
tags with the [development remote](../preferences/dev-remote.md).

## See also

When you run into merge conflicts:

<!-- keep-sorted start -->

- [continue](continue.md) allows you to resume the suspended Git Town command
  after you have resolved the merge conflicts by re-running the failed Git
  command
- [skip](skip.md) ignores all remaining merge conflicts on the current branch
  and then continues the currently suspended Git Town command
- [undo](undo.md) aborts the currently suspended Git Town command and undoes all
  the changes it did, bringing your Git repository back to the state it was
  before you ran the currently suspended Git Town command

<!-- keep-sorted end -->
