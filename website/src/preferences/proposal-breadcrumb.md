# proposal-breadcrumb

This setting controls whether the Git Town CLI embeds a visual representation of
the [branch stack](../stacked-changes.md) (breadcrumbs) into proposals.

If you prefer to handle this outside the CLI, you can achieve the same effect by
using the
[Git Town GitHub Action](https://github.com/marketplace/actions/git-town-github-action).

## values

This setting accepts the following values:

- **none:** do not embed breadcrumbs into proposals
- **branches:** embed breadcrumbs into proposals for all branches
- **stacks:** emded breadcrumbs only for proposals that are part of a stack with
  2 or more branches

## config file

```toml
[sync]
proposal-breadcrumb = "stacks"
```

## Git metadata

To configure whether branches get pushed manually in Git, run this command:

```wrap
git config [--global] git-town.proposal-breadcrumb stacks
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can also configure whether branches get pushed by setting the
`GIT_TOWN_PROPOSAL_BREADCRUMB` environment variable.
