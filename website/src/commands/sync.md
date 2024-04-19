# git sync [--all]

The _sync_ command ("synchronize this branch") updates the local Git workspace
with what happened in the rest of the repository.

- pulls new commits for the current branch from its tracking and ancestor
  branches
- downloads new Git tags
- deletes the local branch if its tracking branch was deleted at the remote and
  the local branch doesn't contain unshipped changes
- checks out the previously checked out Git branch in case the current branch
  got removed as part of the sync
- local branches checked out in other Git worktrees don't get synced

### Arguments

The `--all` parameter makes Git Town sync all local branches instead just the
current one.

The `--dry-run` parameter allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### Configuration

[sync-perennial-strategy](../preferences/sync-perennial-strategy.md) configures
whether perennial branches merge their tracking branch or rebase against it.

[sync-feature-strategy](../preferences/sync-feature-strategy.md) configures
whether feature branches merge their parent and tracking branches or rebase
against them.

If the repository contains a Git remote called `upstream` and the
[sync-upstream](../preferences/sync-upstream.md) setting is enabled, Git Town
also downloads new commits from the upstream main branch.
