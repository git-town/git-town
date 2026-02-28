# proposal-breadcrumb-direction

This setting allows to reverse the direction of the
[breadcrumb](../how-to/proposal-breadcrumb.md) embedded in proposals.

## values

This setting accepts the following values:

- **down:** print the breadcrumb from the root down
- **up:** print the breadcrumb from the root up

## config file

```toml
[propose]
breadcrumb-direction = "down"
```

## Git metadata

To configure whether branches get pushed manually in Git, run this command:

```wrap
git config [--global] git-town.proposal-breadcrumb-direction down
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can also configure whether branches get pushed by setting the
`GIT_TOWN_PROPOSAL_BREADCRUMB_DIRECTION` environment variable.
