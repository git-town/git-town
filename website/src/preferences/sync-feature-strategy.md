# sync-feature-strategy

The sync-feature-strategy setting specifies how to merge changes from their
tracking branch into feature branches.

## values

When set to `merge` (the default value), it merges these changes. When set to
`rebase`, it updates local perennial branches by rebasing them against their
remote branch.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## in config file

In the [config file](../configuration-file.md) the hosting platform is part of
the `[hosting]` section:

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
