# git town sync

> _git town sync [--all] [--stack] [--detached] [--dry-run] [--no-push]_

The _sync_ command ("synchronize this branch") updates the local Git workspace
with what happened in the rest of the repository.

- pulls updates from the tracking and all parent branches
- prunes branches whose tracking branch was deleted at the remote if they
  contain no unshipped changes
- if the branch is checked out in another worktree, syncs with the remote branch
- if the parent branch isn't checked out locally at all, also syncs with its
  remote branch and the parent's parent until it finds a local ancestor branch

Merge conflicts are never fun. If you experience too many merge conflicts, sync
your branches more often. If your Git Town installation is properly configured,
"git town sync --all" syncs all local branches with guarantee to never lose
changes, even in edge cases. You can run it many times per day without thinking
about it. If a sync goes wrong, you can safely go back to the state you repo was
in before the sync by running [git town undo](undo.md).

### --all / -a

By default this command syncs only the current branch. The `--all` aka `-a`
parameter makes Git Town sync all local branches.

### --detached / -d

The `--detached` aka `-d` flag does not pull updates from the main or perennial
branch. This allows you to keep your branches in sync with each other and decide
when to pull in changes from other developers.

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
also downloads new commits from the upstream main branch.

[sync-tags](../preferences/sync-tags.md) configures whether Git Town syncs Git
tags with the `origin` remote.

### Why does git-sync update a branch before deleting it?

"git town sync" can delete branches if their tracking branch was deleted at the
remote. To do so while guaranteeing that it never loses any of your changes,
"git town sync" needs to prove that the branch to be deleted contains no
unshipped changes. Git Town verifies this by diffing the branch to delete
against its parent branch. Only branches with an empty diff can be deleted
safely. For this diff to potentially be empty, Git Town needs to sync the branch
first, even if it's going to be deleted right afterwards.
