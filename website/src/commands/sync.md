# git sync [--all]

The _sync_ command ("synchronize this branch") updates the current branch and
its remote and parent branches with all changes that happened in the repository.
When run on the main or a perennial branch, it pulls and pushes updates and tags
to the tracking branch. When run on a feature branch, it additionally updates
all parent branches and merges the direct parent into the current branch.

If you prefer rebasing your branches instead, set the
[sync-strategy](../preferences/sync-strategy.md) preference.

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../preferences/sync-upstream.md) flag.

### Variations

With the `--all` parameter this command syncs all local branches and not just
the branch you are currently on.

The `--dry-run` parameter allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.
