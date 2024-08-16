# sync-prototype-strategy

The `sync-prototype-strategy` setting specifies how to update local
[prototype branches](../branch-types.md#prototype-branches) with changes from
their parent and tracking branches. When not set, Git Town uses the
[sync-feature-strategy](sync-feature-strategy.md).

## options

`sync-prototype-strategy` accepts the same options as
[sync-feature-strategy](sync-feature-strategy.md#options).

## change this setting

The best way to change this setting is via the
[setup assistant](../configuration.md).

### config file

In the [config file](../configuration-file.md) the sync-prototype-strategy is
part of the `[sync-strategy]` section:

```toml
[sync-strategy]
prototype-branches = "merge"
```

### Git metadata

To manually configure the sync-prototype-strategy in Git, run this command:

```
git config [--global] git-town.sync-prototype-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
