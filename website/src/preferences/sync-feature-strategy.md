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
rebased local commits into the tracking branch. This safe force-push uses the
[--force-with-lease](https://git-scm.com/docs/git-push#Documentation/git-push.txt---no-force-with-lease)
and
[--force-if-includes](https://git-scm.com/docs/git-push#Documentation/git-push.txt---no-force-if-includes)
switches to guarantee that the force-push will never overwrite commits on the
tracking branch that haven't been integrated into your local Git history.

If the safe force-push fails, Git Town rebases your local branch against its
tracking branch. If that leads to conflicts, you have a chance to resolve them
and continue syncing by calling [git town continue](../commands/continue.md).

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
