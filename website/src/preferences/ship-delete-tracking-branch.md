# ship-delete-tracking-branch

```
git-town.ship-delete-tracking-branch=<true|false>
```

If set to `true` (default value), [git ship](../commands/ship.md) deletes the
remote tracking branch of shipped branches.

Some code hosting services like
[GitHub](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-the-automatic-deletion-of-branches)
also delete the tracking branch when shipping via their API. In this case, the
tracking branch is already gone when `git ship` tries to delete it, resulting in
an error. When this happens, change this setting to `false` so that Git Town
does not try to delete the tracking branch.
