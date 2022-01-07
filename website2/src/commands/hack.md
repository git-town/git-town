# Hack command

```
git hack <branch name>
```

The hack command ("start hacking") creates a new feature branch with the given
name off the main branch and brings all uncommitted changes over to it. Before
it does that, it syncs the main branch to ensure the changes in the feature
branch are on top of the latest code version.

If the repository contains a remote called `upstream`, it also syncs the main
branch with its upstream counterpart. You can control this behavior with the
[sync-upstream](../configurations/sync-upstream.md) flag.

If [new-branch-push-flag](.new-branch-push-flag.md) is set, `git hack` creates a
remote tracking branch for the new feature branch. This behavior is disabled by
default to make `git hack` fast. The first run of `git sync` will then create
the remote tracking branch.
