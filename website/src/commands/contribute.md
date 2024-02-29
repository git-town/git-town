# git contribute [branches]

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

Check out a remote branch (that doesn't exists locally on your machine yet) and
make it a contribution branch:

```fish
git contribute other-branch
```

Convert the current contribution branch back to a feature branch:

```fish
git hack
```

Convert the contribution branches "alpha" and "beta" back to a feature branch:

```fish
git hack alpha beta
```
