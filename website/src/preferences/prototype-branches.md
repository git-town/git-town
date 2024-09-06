# prototype-branches

This setting stores the name of all
[prototype branches](../branch-types.md#prototype-branches).

## set prototype branches

To change which branches are prototype in your local repo, use the
[git town park](../commands/park.md) command.

The [Git Town configuration file](../configuration-file.md) doesn't define
define team-wide prototype branches because developers typically only park
branches they own.

## view configured prototype branches

To see which branches are configured as prototype in your local repo, use the
[git town config](../commands/config.md) command.

To check how prototype branches are stored in Git metadata:

```bash
$ git config list | grep 'git-town.prototype-branches'
git-town.prototype-branches=branch-1 branch-2
```
