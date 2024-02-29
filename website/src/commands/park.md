# git park [branches]

The _park_ command [parks](../advanced-syncing#parked-branches) some of your
branches.

## Examples

Park the current branch:

```fish
git park
```

Park branches "alpha" and "beta":

```fish
git park alpha beta
```

Convert the current parked branch back to a feature branch:

```fish
git hack
```

Convert the parked branches "alpha" and "beta" back to feature branches:

```fish
git hack alpha beta
```
