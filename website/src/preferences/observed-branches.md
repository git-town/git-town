# observed-branches

This configuration setting stores the name of all
[observed branches](../branch-types.md#observed-branches).

## set observed branches

The best way to change which branches are observed branches in your local repo
is with the [git town observe](../commands/observe.md) command.

The [Git Town configuration file](../configuration-file.md) cannot and does not
define team-wide observed branches because one developer's feature branch is
another developer's observed or contribution branch.

## view the configured observed branches

The recommended way to see which branches are configured as observed branches in
your local repo is with the [git town config](../commands/config.md) command.

To see how observed branches are stored in Git metadata:

```bash
$ git config list | grep 'git-town.observed-branches'
git-town.observed-branches=branch-1 branch-2
```
