# New branch type

This setting defines the [type](../branch-types.md) for new branches created
using the [git town hack](../commands/hack.md), [append](../commands/append.md),
or [prepend](../commands/prepend.md) commands.

Before setting this, consider one of these more broadly applicable configuration
entries:

- [contribution-regex](contribution-regex.md)
- [unknown-branch-type](unknown-branch-type.md)
- [feature-regex](feature-regex.md)
- [observed-regex](observed-regex.md)
- [perennial-branches](perennial-branches.md)
- [perennial-regex](perennial-regex.md)

## values

These values make sense for this setting:

- `feature` (default)
- `parked`
- `perennial`
- `prototype`

## config file

To configure the type of new branches in the
[configuration file](../configuration-file.md):

```toml
[create]
new-branch-type = "feature"
```

## Git metadata

To configure the type of new branches in Git metadata, run this command:

```wrap
git config [--global] git-town.new-branch-type <feature|parked|perennial|prototype>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the new branch type by setting the `GIT_TOWN_NEW_BRANCH_TYPE`
environment variable.
