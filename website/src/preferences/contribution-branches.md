# contribution-branches

This configuration setting stores the name of all
[contribution branches](../branch-types.md#contribution-branches).

## set contribution branches

The best way to change which branches are contribution branches in your local
repo is with the [git town contribute](../commands/contribute.md) command.

The [Git Town configuration file](../configuration-file.md) cannot and does not
define team-wide contribution branches because one developer's feature branch is
another developer's contribution or observed branch.

## view the configured contribution branches

The recommended way to see which branches are configured as contribution
branches in your local repo is with the [git town config](../commands/config.md)
command.

To see how contribution branches are stored in Git metadata:

```bash
$ git config list | grep 'git-town.contribution-branches'
git-town.contribution-branches=branch-1 branch-2
```
