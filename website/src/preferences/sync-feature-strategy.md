# Feature sync strategy

This setting specifies how to update local feature branches with changes from
their parent and tracking branches.

## options

### merge

When using the "merge" feature sync strategy (which is the default),
[git town sync](../commands/sync.md) merges the parent and tracking branches
into local feature branches.

`merge` is the default value because it is the safest and easiest option.

### rebase

When set to `rebase`, [git town sync](../commands/sync.md) rebases local feature
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
commits. Finish the `git town sync` command and then clean up your commits via a
separate interactive rebase after the sync. At this point another sync will
succeed because the commits you have just cleaned up are now a part of your
local Git history.

The rule of thumb is that pulling in new commits via `git town sync` and
cleaning up old commits must happen separately from each other. Only then can
Git guarantee that the necessary force-push happens without losing commits.

### compress

When using the `compress` sync strategy, [git town sync](../commands/sync.md)
first merges the tracking and parent branches and then
[compresses](../commands/compress.md) the synced branch.

This sync strategy is useful when you want all your pull requests to always
consists of only one commit.

Please be aware that this sync strategy leads to more merge conflicts than the
"merge" sync strategy when more than one Git user makes commits to the same
branch. You can enable these Git settings to prevent this problem:

- `git config rerere.enabled true` enables Git's
  [rerere](https://git-scm.com/book/en/v2/Git-Tools-Rerere) feature
- `git config rerere.autoupdate true` enables auto-staging of auto-resolved
  conflicts

### config file

In the [config file](../configuration-file.md) the feature sync strategy is part
of the `[sync-strategy]` section:

```toml
[sync]
feature-strategy = "merge"
```

### Git metadata

To manually configure the feature sync strategy in Git, run this command:

```wrap
git config [--global] git-town.sync-feature-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
