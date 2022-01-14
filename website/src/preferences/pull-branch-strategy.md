# pull-branch-strategy

```
git-town.pull-branch-strategy [(rebase|merge)]
```

The pull-branch-strategy setting specifies which strategy to use when merging
the remote of the main branch and perennial branches into their local
counterpart. If set to `rebase` (the default value), it updates local perennial
branches via rebase against their remote branch. If set to `merge`, it merges
the respective tracking branch into the local branch.
