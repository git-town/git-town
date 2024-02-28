# git park [branches]

The _park_ command ("park this branch") suspends syncing for the current or
given branches. Possible reasons to do so are:

- you want to intentionally keep the branch at an older state
- you don't want to deal with merge conflicts on this branch right now
- reducing pressure on your CI server

Git Town does not park perennial branches. To unpark, run [hack](hack.md).

## Arguments

```fish
git park
```

Without arguments, Git Town parks the currently checked out branch.
