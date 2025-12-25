# Unknown branch type

This setting defines the branch type to use when Git Town cannot determine the
branch type using all other configuration settings:

- [main-branch](main-branch.md),
- [perennial-branches](perennial-branches.md)
- [Feature regex](feature-regex.md)
- [Contribution regex](contribution-regex.md)
- [Observed regex](observed-regex.md)
- or a manual branch type override set by [git town park](../commands/park.md),
  [git town contribute](../commands/contribute.md),
  [git town hack](../commands/hack.md),
  [git town observe](../commands/observe.md),
  [git town prototype](../commands/prototype.md)

Possible values are:

- `feature` (default)
- `contribution`
- `observed`
- `parked`
- `prototype`

## configuration via setup assistant

A great way to configure this setting is through the setup assistant.

## configure in config file

In the [config file](../configuration-file.md), the unknown branch type is
specified in the `[branches]` section:

```toml
[branches]
unknown-type = "feature"
```

## configure in Git metadata

You can manually configure the unknown branch type using Git metadata:

```wrap
git config [--global] git-town.unknown-branch-type "feature"
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the branch type Git Town should assume for unknown existing
branches by setting the `GIT_TOWN_UNKNOWN_BRANCH_TYPE` environment variable.
