# git compress [branches]

The _contribute_ command makes some of your branches
[contribution](../advanced-syncing.md#contribution-branches) branches.

## Examples

Make the current branch a contribution branch:

```fish
git park
```

Make branches "alpha" and "beta" contribution branches:

```fish
git park alpha beta
```

Check out a remote branch (that exists at origin but not on your local machine)
and make it a contribution branch:

```fish
git contribute somebody-elses-branch
```

Convert the current contribution branch back to a feature branch:

```fish
git hack
```

Convert the contribution branches "alpha" and "beta" back to a feature branch:

```fish
git hack alpha beta
```
