# Sync command

```
git sync [--all]
```

The sync command ("synchronize this branch") updates the current branch and its
remote and parent branches with each other. When run on the main or a perennial
branch, it pulls and pushes updates and tags to the tracking branch. When run on
a feature branch, it additionally updates all parent branches and merges the
direct parent into the current branch.

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../configurations/sync-upstream.md) flag.

The `--all` parameter makes this command sync all local branches and not just
the currently checked out branch. Running with the `--dry-run` parameter prints
the commands but doesn't execute them.
