# sync-perennial-strategy

The sync-perennial-strategy setting specifies how to update local perennial
branches with changes from their tracking branches.

## options

When set to `rebase` (the default value), Git Town rebases local perennial
branches against their tracking branch. When set to `merge`, it merges the
tracking branch into the local perennial branch.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## in config file

In the [config file](../configuration-file.md) the sync-perennial-strategy is
part of the `[sync-strategy]` section:

```toml
[sync-strategy]
perennial-branches = "rebase"
```

## in Git metadata

To manually configure the sync-perennial-strategy in Git, run this command:

```
git config [--global] git-town.sync-perennial-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
