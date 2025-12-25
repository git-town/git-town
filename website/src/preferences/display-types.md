# Display types

This setting allows you to change whether Git Town also displays the
[branch type](../branch-types.md) in addition to the branch name when showing a
list of branches.

## Allowed values

- **all** - display the type for all branches
- **no** - never display the branch type
- **no &lt;branch types&gt;** - display the type of all branches unless the
  branch has one of the listed types
- **&lt;branch types&gt;** - display the type of the branch only if it matches
  one of the listed branch types

## Examples

- `no feature main` displays the type for all branches except for feature and
  main branches. _(This is the default setting.)_
- `prototype observed contribution parked` displays the type only for these four
  branch types

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
your machine. Without it, this setting applies only to the current Git repo.

## Environment variable

You can configure the branch types display via the `GIT_TOWN_SHARE_NEW_BRANCHES`
environment variable.

## CLI flags

You can override this setting per command using:

- `--display-types` / `-d` / `--display-types=all`: display the type for all
  branches
- `--display-types=no`: never display the branch type
- `--display-types=contribution+observed`: display the type only for
  contribution and observed branches
- `--display-types=no-main-feature` displays the type for all branches except
  main and feature branches
