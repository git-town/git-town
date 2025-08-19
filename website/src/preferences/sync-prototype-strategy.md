# Prototype sync strategy

This setting specifies how to update local
[prototype branches](../branch-types.md#prototype-branches) with changes from
their parent and tracking branches. When not set, Git Town uses the
[feature sync strategy](sync-feature-strategy.md).

## options

This setting accepts the same options as the
[feature sync strategy](sync-feature-strategy.md#options).

### config file

In the [config file](../configuration-file.md) the prototype sync strategy is
part of the `[sync-strategy]` section:

```toml
[sync]
prototype-strategy = "merge"
```

### Git metadata

To manually configure the prototype sync strategy in Git, run this command:

```wrap
git config [--global] git-town.sync-prototype-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the sync strategy for prototype branches by setting the
`GIT_TOWN_SYNC_PROTOTYPE_STRATEGY` environment variable.
