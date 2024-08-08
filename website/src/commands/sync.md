# git sync [--all|--stack]

Merge conflicts are never fun, hence minimizing or eliminating them should
always be a priority. To reduce the likelihood of conflicts, it's essential to
keep your branches in sync with each other.

The _sync_ command ("synchronize this branch") updates the local Git workspace
with what happened in the rest of the repository.

- pulls updates from the tracking and all parent branches
- does not modify local branches checked out in other Git worktrees
- deletes branches whose tracking branch was deleted at the remote if they
  contain no unshipped changes

If you experience too many merge conflicts, sync more often. You can run "git
sync" without thinking (and should do so dozens of times per day) because it
guarantees that it never loses any of your changes, even in edge cases. If a
sync goes wrong, you can safely go back to the state you repo was in before the
sync by running [git town undo](undo.md).

### Arguments

By default this command syncs only the current branch. The `--all` parameter
makes Git Town sync all local branches. The `--stack` parameter makes Git Town
sync all branches in the stack that the current branch belongs to.

The `--dry-run` parameter allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

The `--no-push` argument disables all pushes of local commits to their tracking
branch.

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

"git sync" can delete branches if their tracking branch was deleted at the
remote. To do so while guaranteeing that it never loses any of your changes,
"git sync" needs to prove that the branch to be deleted contains no unshipped
changes. Git Town verifies this by diffing the branch to delete against its
parent branch. Only branches with an empty diff can be deleted safely. For this
diff to potentially be empty, Git Town needs to sync the branch first, even if
it's going to be deleted right afterwards.
