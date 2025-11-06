# Share new branches

This setting allows you to change how Git Town shares new branches created with
[hack](../commands/hack.md), [append](../commands/append.md), or
[prepend](../commands/prepend.md).

Allowed values:

- **no/false/0:** Don't share new branches, keep them local on your machine
  until you [sync](../commands/sync.md) or [propose](../commands/propose.md)
  them (default behavior).
- **push:** Push new branches to the [development remote](dev-remote.md).
- **propose:** Create a pull request for the new branch. This is similar to
  always adding the [propose flag](../commands/hack.md#--propose).

## in config file

```toml
create.share-new-branches = "push|propose"
```

## in Git metadata

To enable pushing new branches in Git, run this command:

```wrap
git config [--global] git-town.share-new-branches <push|propose>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.

## environment variable

You can configure how new branches get shared by setting the
`GIT_TOWN_SHARE_NEW_BRANCHES` environment variable.
