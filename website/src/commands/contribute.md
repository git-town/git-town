# git town contribute

> _git town contribute [branch-name...]_

The _contribute_ command makes some of your branches
[contribution](../branch-types.md#contribution-branches) branches.

When called without arguments, it makes the current branch a contribution
branch.

To convert a contribution branch back into a feature branch, use the
[hack](hack.md) command.

To make the current branch a contribution branch:

```fish
git contribute
```

### Positional arguments

When called with positional arguments, this commands makes the branches with the
given names contribution branches.

To make branches "alpha" and "beta" contribution branches:

```fish
git contribute alpha beta
```

Check out a remote branch (that exists at origin but not on your local machine)
and make it a contribution branch:

```fish
git contribute somebody-elses-branch
```

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
