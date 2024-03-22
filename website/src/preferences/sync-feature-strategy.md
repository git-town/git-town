# sync-feature-strategy

The `sync-feature-strategy` setting specifies how to update local feature
branches with changes from their parent and tracking branches.

## options

### merge

When using the "merge" sync-feature-strategy, [git sync](../commands/sync.md)
merges the parent and tracking branches into local feature branches.

`merge` is the default value because it is the safest and easiest option.

### rebase

When set to `rebase`, [git sync](../commands/sync.md) rebases local feature
branches against their parent branches and then does a safe force-push of your
rebased local commits to the tracking branch. This safe force-push uses Git's
[--force-with-lease](https://git-scm.com/docs/git-push#Documentation/git-push.txt---no-force-with-lease)
and
[--force-if-includes](https://git-scm.com/docs/git-push#Documentation/git-push.txt---no-force-if-includes)
switches to guarantee that the force-push will never overwrite commits on the
tracking branch that haven't been integrated into the local Git history.

If the safe force-push fails, Git Town rebases your local branch against its
tracking branch to pull in new commits from the tracking branch. If that leads
to conflicts, you have a chance to resolve them and continue syncing by running
[git town continue](../commands/continue.md).

When continuing the sync this way, Git Town tries again to safe-force-push and
rebase until the safe-force-push succeeds without removing commits from the
tracking branch that aren't part of the local Git history.

This can lead to an infinite loop if you do an interactive rebase that removes
commits from the tracking branch while syncing it. You can break out of this
infinite loop by doing a less aggressive rebase that doesn't remove the remote
commits. Finish the `git sync` command and then clean up your commits via a
separate interactive rebase after the sync. At this point another sync will
succeed because the commits you have just cleaned up are now a part of your
local Git history.

The rule of thumb is that pulling in new commits and cleaning up old commits
must happen separately from each other. Only then can Git guarantee that pulling
in changes happens without losing commits.

## change this setting

The best way to change this setting is via the
[setup assistant](../configuration.md).

### config file

In the [config file](../configuration-file.md) the sync-feature-strategy is part
of the `[sync-strategy]` section:

```toml
[sync-strategy]
feature-branches = "merge"
```

### Git metadata

To manually configure the sync-feature-strategy in Git, run this command:

```
git config [--global] git-town.sync-feature-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
