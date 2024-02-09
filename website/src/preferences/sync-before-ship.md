# sync-before-ship

Syncing before shipping allows you to resolve merge conflicts and any associated
problems on the feature branch instead of the main branch. This helps keep the
main branch green. The downside of syncing before shipping is that it will
trigger another CI run and might block shipping until CI is green again. Syncing
before shipping therefore makes the most sense when shipping locally on your
machine.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## values

When set to `true`, `git ship` syncs the branch to ship before shipping. When
set to `false` (the default value), `git ship` does not sync the branch before
shipping.

## in config file

To configure `sync-before-ship` in the
[configuration file](../configuration-file.md):

```toml
sync-before-ship = false
```

## in Git metadata

To manually configure `sync-before-ship` in Git, run this command:

```
git config [--global] git-town.sync-before-ship <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
