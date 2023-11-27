# ship-delete-remote-branch

```
git-town.ship-delete-remote-branch=<true|false>
```

If set to `true` (default value), [git ship](../commands/ship.md) deletes the
remote tracking branch of shipped branches. Some code hosting services like
[GitHub](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-the-automatic-deletion-of-branches)
also delete the remote branch when merging a proposal. In this case, change this
setting to `false` so that Git Town skips deleting the tracking branch.
