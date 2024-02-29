# git contribute [branches]

The _contribute_ command makes some of your branches contribution branches.
Parked branches don't get synced unless you run `git sync` directly on a parked
branch.

You might want to park a branch if you

- want to intentionally keep the branch at an older state
- don't want to deal with merge conflicts on this branch right now
- reduce load on your CI server

Perennial branches and the main branch cannot get parked. Use
[git hack](hack.md) to unpark branches ("I want to resume hacking on this
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
