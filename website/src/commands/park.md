# git park [branches]

The _park_ command parks some of your branches. Parked branches don't get synced
unless you run `git sync` directly on a parked branch.

You might want to park a branch if you

- want to intentionally keep the branch at an older state
- don't want to deal with merge conflicts on this branch right now
- reduce load on your CI server

Perennial branches and the main branch cannot get parked. Use
[git hack](hack.md) to unpark branches.

## Arguments

```fish
git park
```

Without arguments, Git Town parks the currently checked out branch.
