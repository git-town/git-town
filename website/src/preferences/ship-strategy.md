# ship-strategy

The `ship-strategy` setting defines how [git town ship](../commands/ship.md)
merges finished feature branches into the main branch.

## options

### api

When using the "api" ship-strategy, [git ship](../commands/ship.md) presses the
"merge" button for the proposal in the web UI of your code hosting platform via
an API call.

You need to configure an API token in the
[setup assistant](../commands/config-setup.md) for this to work.

`api` is the default value because it does exactly what you normally do
manually.

### squash-merge

When set to `squash-merge`, [git ship](../commands/ship.md) merges the feature
branch to ship in your local Git repository. While doing so it squashes all
commits on the feature branch into a single commit and lets you edit the commit
message.

## change this setting

The best way to change this setting is via the
[setup assistant](../configuration.md).

### config file

Set the ship-strategy in the [config file](../configuration-file.md):

```toml
ship-strategy = "api"
```

### Git metadata

To manually configure the ship-strategy in Git metadata, run this command:

```
git config [--global] git-town.ship-strategy <api|squash-merge>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
