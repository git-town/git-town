# Perennial sync strategy

This setting specifies how to update local perennial branches with changes from
their tracking branches.

## options

### rebase

When using the `rebase` sync strategy, (which is the default), Git Town rebases
local perennial branches onto their tracking branch.

### ff-only

Git Town fast-forwards the local branch to match the tracking branch. If a
fast-forward is not possible, Git Town exits with a descriptive error message.
This is ideal when you want an explicit warning about unpushed local commits.

## in config file

In the [config file](../configuration-file.md) the perennial sync strategy is
part of the `[sync-strategy]` section:

```toml
[sync]
perennial-strategy = "rebase"
```

## in Git metadata

To manually configure the perennial sync strategy in Git, run this command:

```wrap
git config [--global] git-town.sync-perennial-strategy <ff-only|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the sync strategy for perennial branches by setting the
`GIT_TOWN_SYNC_PERENNIAL_STRATEGY` environment variable.
