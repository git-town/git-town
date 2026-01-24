# proposal-breadcrumb

This setting makes the Git Town CLI embed a visual overview of the
[branch stack](../stacked-changes.md) into proposals.

An alternative to doing this is setting up the
[Git Town GitHub action](https://github.com/marketplace/actions/git-town-github-action).

## config file

```toml
[sync]
proposal-breadcrumb = "cli"
```

## Git metadata

To configure whether branches get pushed manually in Git, run this command:

```wrap
git config [--global] git-town.proposal-breadcrumb cli
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure whether branches get pushed by setting the
`GIT_TOWN_PROPOSAL_BREADCRUMB` environment variable.
