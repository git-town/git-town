# git contribute [branches]

The _contribute_ command makes some of your branches
[contribution](../advanced-syncing#contribution-branches) branches.

## Examples

Make the current branch a contribution branch:

```fish
git park
```

Make branches "alpha" and "beta" contribution branches:

```fish
git park alpha beta
```

Convert the current contribution branch back to a feature branch:

```fish
git hack
```

Convert the contribution branches "alpha" and "beta" back to a feature branch:

```fish
git hack alpha beta
```
