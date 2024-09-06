# parked-branches

This configuration setting stores the name of all
[parked branches](../branch-types.md#parked-branches).

## set parked branches

The best way to change which branches are parked branches in your local repo is
with the [git town park](../commands/park.md) command.

The [Git Town configuration file](../configuration-file.md) cannot and does not
define team-wide parked branches because typically developers park only branches
which they own.

## view the configured parked branches

The recommended way to see which branches are configured as parked branches in
your local repo is with the [git town config](../commands/config.md) command.

To see how parked branches are stored in Git metadata:

```bash
$ git config list | grep 'git-town.parked-branches'
git-town.parked-branches=branch-1 branch-2
```
