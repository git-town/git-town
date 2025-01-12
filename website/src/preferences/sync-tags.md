# Sync tags

This setting configures whether to sync Git tags with the
[development remote](dev-remote.md).

## options

When set to `true` (the default value), `git town sync` also pulls and pushes
Git tags in addition to branches and commits. When set to `false`,
`git town sync` does not change Git tags at the local or remote Git repository.

## in config file

In the [config file](../configuration-file.md) syncing tags can be set like
this:

```toml
[sync]
tags = true
```

## in Git metadata

To manually configure syncing tags in Git, run this command:

```wrap
git config [--global] git-town.sync-tags <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
