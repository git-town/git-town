# new-branch-type configuration setting

The "new-branch-type" setting defines the [type](../branch-types.md) for new
branches created using the [git town hack](../commands/hack.md),
[append](../commands/append.md), or [prepend](../commands/prepend.md) commands.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## values

- `feature` (default)
- `parked`
- `perennial`
- `prototype`

## config file

To configure the type of new branches in the
[configuration file](../configuration-file.md):

```toml
new-branch-type = "feature"
```

## Git metadata

To configure the type of new branches in Git metadata, run this command:

```bash
git config [--global] git-town.new-branch-type <feature|parked|perennial|prototype>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
