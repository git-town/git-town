# Delete tracking branch

Some forges like
[GitHub](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-the-automatic-deletion-of-branches)
and
[GitLab](http://ncugw.phy.ncu.edu.tw/gitlab/help/user/project/merge_requests/getting_started.md#deleting-the-source-branch)
can delete the tracking branch when shipping via their API. In this case the
tracking branch is already gone when `git town ship` tries to delete it,
resulting in an error. To prevent this error, set this setting to `false` so
that Git Town does not try to delete the tracking branch.

## in config file

```toml
ship.delete-tracking-branch = true
```

or

```toml
[ship]
delete-tracking-branch = true
```

## in Git metadata

To configure this setting in Git, run this command:

```wrap
git config [--global] git-town.ship-delete-tracking-branch <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
