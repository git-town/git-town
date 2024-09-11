# observed-branches

This configuration setting stores the name of all
[observed branches](../branch-types.md#observed-branches).

## set observed branches

To change which branches are observed branches in your local repo, use the
[git town contribute](../commands/contribute.md) command.

The [Git Town configuration file](../configuration-file.md) does not define
team-wide observed branches because one developer's feature branch is another
developer's contribution or observed branch.

## view configured observed branches

To see which branches are configured as observed branches in your local repo,
use the [git town config](../commands/config.md) command.

To view how observed branches are stored in Git metadata:

```bash
$ git config list | grep 'git-town.observed-branches'
git-town.observed-branches=branch-1 branch-2
```
