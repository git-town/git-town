# ignore-uncommitted

By default, Git Town refuses to [ship](../commands/ship.md) a branch if there
are uncommitted changes, ensuring that everything on the branch is included in
the ship. This setting allows you to configure this behavior.

## options

- `false` (default) requires a clean workspace. This guarantees that all changes
  on the branch are committed and shipped.
- `true` allows shipping with uncommitted changes, i.e. what CI sees.

## via CLI flag

You can override the configured behavior for a single invocation:

```sh
git-town ship --ignore-uncommitted
git-town ship --no-ignore-uncommitted
```

## in config file

To configure this behavior permanently, you can configure it in the
[config file](../configuration-file.md):

```toml
[ship]
ignore-uncommitted = true
```

## in Git metadata

You can also configure this setting via Git config:

```wrap
git config [--global] git-town.ignore-uncommitted <true|false>
```

The optional `--global` flag applies this setting to all repositories on your
machine. Without it, the setting applies only to the current repository.

## environment variable

You can control this behavior by setting the `GIT_TOWN_IGNORE_UNCOMMITTED`
environment variable.

```sh
env GIT_TOWN_IGNORE_UNCOMMITTED=true git-town ship
```
