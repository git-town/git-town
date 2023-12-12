# git sync [--all]

The _sync_ command ("synchronize this branch") updates the currently checked out
branch and its remote and parent branches with all changes that happened in the
repository.

When run on the main or a perennial branch, it pulls and pushes updates and tags
to the tracking branch. When run on a feature branch, it additionally syncs all
parent branches and merges/rebases the direct parent into the current branch. If
the branch was deleted at the remote, and the local branch contains no unshipped
changes, Git Town removes it from the local workspace.

### Arguments

With the `--all` parameter this command syncs all local branches and not just
the current one.

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
