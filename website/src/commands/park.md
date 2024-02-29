# git park [branches]

The _park_ command [parks](../advanced-syncing#parked-branches) some of your
branches.

Use [git hack](hack.md) to unpark branches ("I want to continue hacking on this
branch").

## Examples

Park the current branch:

```fish
git park
```

Park branches "alpha" and "beta":

```fish
git park alpha beta
```

Unpark the current branch:

```fish
git hack
```

Unpark branches "alpha" and "beta":

```fish
git hack alpha beta
```
