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

### fast-forward

When using the `fast-forward` ship strategy, the `git ship` command merges your
feature branch into the parent branch using the `--ff-only` option. This ensures
that the parent branch will have the exact same commits as your feature branch.

The `fast-forward` strategy is ideal for shipping
[stacked changes](../stacked-changes.md) because it avoids the empty merge
conflicts that can occur when squash-merging in this context. For more details,
check out this
[GitHub documentation](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges#squashing-and-merging-a-long-running-branch).

Performing a local merge and pushing the new commits on your main branch works
on GitHub even if your main branch is protected! GitHub recognizes that the
commits have already been tested by CI and allows them to be pushed. Just make
sure the branch passes CI tests. For more information, see
[this StackOverflow answer](https://stackoverflow.com/questions/60597400/how-to-do-a-fast-forward-merge-on-github/66906599#66906599).

A limitation of the `fast-forward` ship strategy is that your feature branch
must be up to date, i.e. the main branch must not have received additional
commits since you last synced your feature branch.

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

To manually configure the ship-strategy in Git metadata, run:

```
git config [--global] git-town.ship-strategy <api|squash-merge>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
