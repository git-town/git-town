# parked-branches

This setting stores the name of all
[parked branches](../branch-types.md#parked-branches).

## set parked branches

To change which branches are parked in your local repo, use the
[git town park](../commands/park.md) command.

The [Git Town configuration file](../configuration-file.md) doesn't define
team-wide parked branches because developers typically only park branches they
own.

## view configured parked branches

To see which branches are configured as parked in your local repo, use the
[git town config](../commands/config.md) command.

To view how parked branches are stored in Git metadata:

```bash
$ git config list | grep 'git-town.parked-branches'
git-town.parked-branches=branch-1 branch-2
```
