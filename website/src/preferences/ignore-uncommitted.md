# ignore-uncommitted

This setting configures whether Git Town requires no committed files when
[shipping branches](../commands/ship.md).

## options

When set to `false` (the default value), [git town ship](../commands/ship.md)
requires no uncommitted changes. This ensures all changes on that branch get
shipped. When set to`true`, you can ship with uncommitted changes.

## via CLI flag

```sh
git-town ship --ignore-uncommitted
```

## in config file

The [config file](../configuration-file.md) can enable detached mode permanently
for all commands like this:

```toml
[ship]
ignore-uncommitted = true
```

## in Git metadata

```wrap
git config [--global] git-town.ignore-uncommitted <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure whether Git Town syncs Git tags by setting the
`GIT_TOWN_IGNORE_UNCOMMITTED` environment variable.
