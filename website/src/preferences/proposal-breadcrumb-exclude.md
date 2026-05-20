# Proposal-breadcrumb-exclude

Controls which branch types Git Town omits from proposal breadcrumbs.

This setting only affects breadcrumbs
when [proposal-breadcrumb](proposal-breadcrumb.md) is set to `branches`
or `stacks`.

## options

Supported values are defined in [branch types](../branch-types.md).

An empty list excludes no branch types.

## config file

To configure this in the configuration file, define the excluded branch types:

```toml
breadcrumb-exclude-branches = ["contribution", "prototype"]
```

## Git metadata

To configure this in Git, run this command:

```wrap
git config [--global] git-town.proposal-breadcrumb-exclude "contribution, prototype"
```

The optional `--global` flag applies this setting to all Git repositories on
your machine.
Without it, the setting applies only to the current repository.

## environment variable

You can also configure this via the environment variable
`GIT_TOWN_PROPOSAL_BREADCRUMB_EXCLUDE`

```sh
export GIT_TOWN_PROPOSAL_BREADCRUMB_EXCLUDE="contribution prototype"
```
