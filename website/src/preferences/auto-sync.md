# auto-sync

By default, Git Town automatically keeps your branches in sync with each other,
for example before creating a new branch. This setting allows you to disable
this behavior. When disabled, your branches are no longer automatically synced
and you need to run [git town sync](../commands/sync.md) manually to sync them.

## CLI flags

In one-off situations you can call commands that sync branches with the
`--no-sync` flag to disable automatic syncing.

If you have automatic syncing disabled permanently via the config file or Git
metadata (see below), you can enable it with the `--sync` flag.

## config file

To configure automatic syncing in the
[configuration file](../configuration-file.md):

```toml
[sync]
auto-sync = false
```

## Git metadata

To configure automatic syncing in Git metadata, run this command:

```wrap
git config [--global] git-town.auto-sync <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure auto-sync by setting the `GIT_TOWN_AUTO_SYNC` environment
variable.
