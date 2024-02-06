# sync-upstream

If you The sync-upstream setting configures whether to pull in updates from the
`upstream` remote. This applies to codebases that are forks of other codebases.

## options

When set to `true` (the default value), `git sync` also updates the local
[main main-branch](main-branch.md) with changes from its counterpart in the
`upstream` remote. When set to `false`, `git sync` does not pull in updates from
upstream even if that remote exists.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## in config file

In the [config file](../configuration-file.md) the sync-upstream setting can be
set like this:

```toml
sync-upstream = true
```

## in Git metadata

To manually configure `sync-upstream` in Git, run this command:

```
git config [--global] git-town.sync-upstream <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
