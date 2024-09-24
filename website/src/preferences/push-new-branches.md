# push-new-branches

By default, Git Town does not push new feature branches to the `origin` remote
since that would make creating branches slower and triggers an unnecessary CI
run for a branch containing no changes. Running
[git town sync](../commands/sync.md) or
[git town propose](../commands/propose.md) later will push the branch to origin.
If you prefer to push new branches upon creation, enable this configuration
option.

## in config file

```toml
push-new-branches = true
```

## in Git metadata

To enable pushing new branches in Git, run this command:

```
git config [--global] push-new-branches <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.
