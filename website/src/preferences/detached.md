# Detached

This setting configures whether Git Town pulls new commits from the main or
perennial branch at the root of your branch hierarchy. This can help if you
encounter too many interruptions through expensive recompiles in busy monorepos
after syncing.

## options

When set to `false` (the default value), `git town sync` pulls updates from the
perennial root branch of your stack. When set to `true`, `git town sync` does
not pull in changes from the perennial root.

## via CLI flag

Commands that sync branches have a `--detached` flag that your can set to
perform that sync in detached mode. If you have enabled detached mode
permenantly through the configuration setting described on this page, you can
use the `--no-detached` CLI flag to disable detached mode for that command.

## in config file

In the [config file](../configuration-file.md) syncing tags can be set like
this:

```toml
[sync]
detached = true
```

## in Git metadata

To manually configure syncing tags in Git, run this command:

```wrap
git config [--global] git-town.detached <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure whether Git Town syncs Git tags by setting the
`GIT_TOWN_DETACHED` environment variable.
