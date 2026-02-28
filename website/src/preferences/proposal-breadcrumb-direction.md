# proposal-breadcrumb-direction

Controls the direction of the [breadcrumb](../how-to/proposal-breadcrumb.md)
that Git Town embeds into proposals.

## values

Supported values:

- **down** renders the breadcrumb from the root branch down
- **up** renders the breadcrumb from the root branch up

Choose the direction that best matches how your team thinks about stacked
branches.

## config file

```toml
[propose]
breadcrumb-direction = "down"
```

## Git metadata

To configure this via Git metadata:

```wrap
git config [--global] git-town.proposal-breadcrumb-direction down
```

With `--global`, the setting applies to all Git repositories on your machine.
Without it, the setting applies only to the current repository.

## environment variable

You can also configure this via the environment variable
`GIT_TOWN_PROPOSAL_BREADCRUMB_DIRECTION`.
