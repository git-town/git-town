# git observe [branches]

The _observe_ command makes some of your branches
[observed](../advanced-syncing.md#observed-branches) branches.

## Examples

Observe the current branch:

```fish
git observe
```

Observe branches "alpha" and "beta":

```fish
git observe alpha beta
```

Check out a remote branch (that doesn't exists locally on your machine yet) and
make it observed:

```fish
git observe other-branch
```

Convert the current observed branch back to a feature branch:

```fish
git hack
```

Convert the observed branches "alpha" and "beta" back to feature branches:

```fish
git hack alpha beta
```
