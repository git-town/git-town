# git town sync

> _git town sync [--all] [--detached] [--dry-run] [--no-push] [--stack]
> [--verbose]_

The _sync_ command ("synchronize this branch") updates your local Git workspace
with what happened in the rest of the repository.

Merge conflicts are not fun and can break code. Minimize them by syncing your
branches frequently. Git town knows how to sync many different types of
branches. When properly configured, `git town sync --all` can synchronize all
your local branches the right way without losing changes, even in edge cases.

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

### --all / -a

By default this command syncs only the current branch. The `--all` aka `-a`
parameter makes Git Town sync all local branches.

### --detached / -d

The `--detached` aka `-d` flag does not pull updates from the main or perennial
branch at the root of your branch hierarchy. This allows you to keep your
branches in sync with each other and decide when to pull in changes from other
developers.

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --no-push

The `--no-push` argument disables all pushes of local commits to their tracking
branch.

### --stack / -s

The `--stack` aka `-s` parameter makes Git Town sync all branches in the stack
that the current branch belongs to.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

### Configuration

[sync-perennial-strategy](../preferences/sync-perennial-strategy.md) configures
whether perennial branches merge their tracking branch or rebase against it.

[sync-feature-strategy](../preferences/sync-feature-strategy.md) configures
whether feature branches merge their parent and tracking branches or rebase
against them.

If the repository contains a Git remote called `upstream` and the
[sync-upstream](../preferences/sync-upstream.md) setting is enabled, Git Town
also pulls new commits from the upstream's main branch.

[sync-tags](../preferences/sync-tags.md) configures whether Git Town syncs Git
tags with the `origin` remote.

### Why does git-sync sometimes update a local branch whose tracking branch was deleted before deleting it?

If a remote branch was deleted at the remote, it is considered obsolete and "git
town sync" will remove its local counterpart. To guarantee that this doesn't
lose unshipped changes in the local branch, "git town sync" needs to prove that
the branch to be deleted contains no unshipped changes.

The easiest way to prove that is when the local branch was in sync with its
tracking branch before Git Town runs `git fetch`. This is another reason to run
`git town sync` regularly.

If a local shipped branch is not in sync with its tracking branch on your
machine, Git Town must check for unshipped local changes by diffing the branch
to delete against its parent branch. Only branches with an empty diff can be
deleted safely. For this to work, Git Town needs to sync the branch first, even
if it's going to be deleted right afterwards.
