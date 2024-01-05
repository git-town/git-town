# ship-delete-tracking-branch

```
git-town.ship-delete-tracking-branch=<true|false>
```

If set to `true` (default value), [git ship](../commands/ship.md) deletes the
remote tracking branch of shipped branches.

Some code hosting services like
[GitHub](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-the-automatic-deletion-of-branches)
and
[GitLab](http://ncugw.phy.ncu.edu.tw/gitlab/help/user/project/merge_requests/getting_started.md#deleting-the-source-branch)
can also delete the tracking branch when shipping via their API. In this case
the tracking branch is already gone when `git ship` tries to delete it,
resulting in an error. To prevent this error, change the
_ship-delete-tracking-branch_ setting to `false` so that Git Town does not try
to delete the tracking branch.
