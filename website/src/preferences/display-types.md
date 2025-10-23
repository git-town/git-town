# Display types

This setting allows you to change whether and how Git Town displays
[branch types](../branch-types.md).

## Allowed values

- **no** - display no branch types
- **all** - display all branch types
- **no &lt;branch types&gt;** - display all branch types except the given ones
- **&lt;branch types&gt;** - display only the given branch types

## Examples

- `no feature main` displays all branch types except if the branch is a feature
  branch or the main branch. This is the default setting.
- `prototype observed contribution parked` displays only the given four branch
  types

## Config file

```toml
[branches]
display-types = "<value>"
```

## Git metadata

```wrap
git config [--global] git-town.display-types <push|propose>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, this setting applies to the current Git repo.

## Environment variable

You can configure which branch types Git Town displays via the
`GIT_TOWN_SHARE_NEW_BRANCHES` environment variable.

## CLI flags

You can override this setting per command using:

- `--display-types` / `-d` / `--display-types=all` displays all types
- `--display-types=no` displays no types
- `--display-types=contribution+observed` displays only the types for
  contribution and observed branches
- `--display-types=no-main-feature` displays all branch types except for the
  main and feature branches
