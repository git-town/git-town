# New branch type

This setting defines the [type](../branch-types.md) for new branches created
using the [git town hack](../commands/hack.md), [append](../commands/append.md),
or [prepend](../commands/prepend.md) commands.

Before setting this, try to configure the branch type using one of these more
broadly applicable configuration entries:

- [contribution-regex](contribution-regex.md)
- [default-branch-type](default-branch-type.md)
- [feature-regex](feature-regex.md)
- [observed-regex](observed-regex.md)
- [perennial-branches](perennial-branches.md)
- [perennial-regex](perennial-regex.md)

## values

- `feature` (default)
- `parked`
- `perennial`
- `prototype`

## config file

To configure the type of new branches in the
[configuration file](../configuration-file.md):

```toml
new-branch-type = "feature"
```

## Git metadata

To configure the type of new branches in Git metadata, run this command:

```wrap
git config [--global] git-town.new-branch-type <feature|parked|perennial|prototype>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
