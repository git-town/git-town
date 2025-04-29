# Share new branches

By default, Git Town doesn't push new feature branches to the
[development remote](dev-remote.md) since that would make creating branches
slower and triggers an unnecessary CI run for a branch containing no changes.
Running [git town sync](../commands/sync.md) or
[git town propose](../commands/propose.md) later will push the branch to the dev
remote. If you prefer to push new branches upon creation, set this configuration
option.

## in config file

```toml
create.share-new-branches = "push"
```

## in Git metadata

To enable pushing new branches in Git, run this command:

```wrap
git config [--global] share-new-branches push
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.
