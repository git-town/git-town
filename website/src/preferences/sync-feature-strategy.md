# sync-feature-strategy

The `sync-feature-strategy` setting specifies how to update local feature
branches with changes from their parent and tracking branches.

## options

When set to `merge` (the default value), Git Town merges the parent and tracking
branches into local feature branches. When set to `rebase`, it rebases local
feature branches against their parent and tracking branches.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## in config file

In the [config file](../configuration-file.md) the sync-feature-strategy is part
of the `[sync-strategy]` section:

```toml
[sync-strategy]
feature-branches = "merge"
```

## in Git metadata

To manually configure the sync-feature-strategy in Git, run this command:

```
git config [--global] git-town.sync-feature-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
