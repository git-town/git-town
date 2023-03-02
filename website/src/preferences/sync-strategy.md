# sync-strategy

```
git-town.sync-strategy <merge|rebase>
```

The sync-strategy setting specifies which strategy to use when merging the
remote of feature branches into their local counterpart. If set to `merge` (the
default value), it merges the respective tracking branch into its local branch.
If set to `rebase`, it updates local perennial branches by rebasing them against
their remote branch.
