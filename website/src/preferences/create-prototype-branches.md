# create-prototype-branches configuration setting

The "create-prototype-branches" setting determines whether Git Town creates new
branches as [prototype branches](../branch-types.md#prototype-branches).

The best way to change this setting is via the
[setup assistant](../configuration.md).

## config file

To configure the creation of prototype branches in the
[configuration file](../configuration-file.md):

```toml
create-prototype-branches = false
```

## Git metadata

To configure the creation of prototype branches manually in Git, run this
command:

```bash
git config [--global] git-town.create-prototype-branches <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
