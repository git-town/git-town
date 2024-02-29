# git observe [branches]

The _observe_ command makes some of your branches
[observed](../advanced-syncing#observed-branches) branches.

## Examples

Observe the current branch:

```fish
git observe
```

Observe branches "alpha" and "beta":

```fish
git observe alpha beta
```

Un-observe the current branch (makes it an
[owned branch](../advanced-syncing#branch-ownership)):

```fish
git hack
```

Un-observe branches "alpha" and "beta":

```fish
git hack alpha beta
```
