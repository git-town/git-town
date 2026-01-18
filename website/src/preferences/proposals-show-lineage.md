# proposals-show-lineage

This setting configures how Git Town embeds a
[breadcrumb](../how-to/github-actions-breadcrumb.md) of the
[branch stack](../stacked-changes.md) into proposals.

You have several options for this:

1. Let the Git Town executable create and maintain branch lineages of proposals.
2. Use the
   [Git Town GitHub action](https://github.com/marketplace/actions/git-town-github-action)

## config file

```toml
[sync]
proposals-show-lineage = cli
```

## Git metadata

To configure whether branches get pushed manually in Git, run this command:

```wrap
git config [--global] git-town.proposals-show-lineage cli
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure whether branches get pushed by setting the
`GIT_TOWN_PROPOSALS_SHOW_LINEAGE` environment variable.
