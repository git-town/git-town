# Ship a branch from a stack without re-running CI for all

When you have a [stack](../stacked-changes.md) of branches, and you ship the
oldest branch, normally you have to run `git sync --stack` to sync the just
shipped changes through the stack before you can ship the next branch in the
stack. This forces a re-run of CI, which is wasteful because no code was
changed.

You can avoid re-running CI tests by shipping via a
[fast-forward merge](../preferences/ship-strategy.md#fast-forward):

```
git switch <branch to ship>
git ship --strategy=fast-forward
```

This requires that the main branch has received no updates since the stack was
synced.
