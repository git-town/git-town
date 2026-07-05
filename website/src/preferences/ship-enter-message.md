# ship-enter-message

When [ship-strategy](ship-strategy.md) is set to `api`,
`git town ship` merges the proposal for the current branch via your forge's API.
By default it lets the forge determine the squash commit message,
the same way clicking the "merge" button in the forge's web UI does,
without any input from you.

Enable this setting to enter the squash commit message yourself instead.
Git Town then performs a local squash merge
so you can edit the commit message in your editor,
and uses that message when merging the proposal via the API.

You can also enter the message
for a single ship with the
[--enter-message](../commands/ship.md#--enter-message) flag, or provide it
directly with [--message](../commands/ship.md#-m-textmessage-text).

## config file

```toml
ship.enter-message = true
```

or

```toml
[ship]
enter-message = true
```

## Git metadata

To configure this setting in Git, run this command:

```wrap
git config [--global] git-town.ship-enter-message <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine.
Without it, the setting applies only to the current repository.

## environment variable

You can configure whether ship lets you enter the commit message by setting the
`GIT_TOWN_SHIP_ENTER_MESSAGE` environment variable.
