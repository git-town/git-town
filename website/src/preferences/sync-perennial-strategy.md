# sync-perennial-strategy

The sync-perennial-strategy setting specifies how to update local perennial
branches with changes from their tracking branches.

## options

- `rebase` (default value): Git Town rebases local perennial branches against
  their tracking branch.
- `ff-only`: it fast-forwards the local branch to match the tracking branch and
  exits with an error messages if the local branch contains unpushed commits
- `merge`: it merges the tracking branch into the local perennial branch

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
