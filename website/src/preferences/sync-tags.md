# sync-tags

The sync-tags setting configures whether to sync Git tags with the `origin`
remote.

## options

When set to `true` (the default value), `git sync` also pulls and pushes Git
tags in addition to branches and commits. When set to `false`, `git sync` does
not change Git tags at the local or remote Git repository.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## in config file

In the [config file](../configuration-file.md) the sync-tags setting can be set
like this:

```toml
sync-tags = true
```

## in Git metadata

To manually configure `sync-tags` in Git, run this command:

```
git config [--global] git-town.sync-tags <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
