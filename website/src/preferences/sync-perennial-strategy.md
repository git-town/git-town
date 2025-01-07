# sync-perennial-strategy

The sync-perennial-strategy setting specifies how to update local perennial
branches with changes from their tracking branches.

## options

- `rebase` (default value): Git Town rebases local perennial branches against
  their tracking branch.
- `ff-only`: Git Town fast-forwards the local branch to match the tracking
  branch. If a fast-forward is not possible, it exits with an error message.
  This strategy is great if you want get warned about unpushed local commits.
- `merge`: Git Town merges the tracking branch into the local perennial branch.

## in config file

In the [config file](../configuration-file.md) the sync-perennial-strategy is
part of the `[sync-strategy]` section:

```toml
[sync]
perennial-strategy = "rebase"
```

## in Git metadata

To manually configure the sync-perennial-strategy in Git, run this command:

```wrap
git config [--global] git-town.sync-perennial-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
