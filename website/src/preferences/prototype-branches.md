# prototype-branches

This configuration setting stores the name of all
[prototype branches](../branch-types.md#prototype-branches).

## set prototype branches

The best way to change which branches are prototype branches in your local repo
is with the [git town prototype](../commands/prototype.md) command.

The [Git Town configuration file](../configuration-file.md) cannot and does not
define team-wide prototype branches because typically developers prototype only
branches which they own.

## view the configured prototype branches

The recommended way to see which branches are configured as prototype branches
in your local repo is with the [git town config](../commands/config.md) command.

To see how prototype branches are stored in Git metadata:

```bash
$ git config list | grep 'git-town.prototype-branches'
git-town.prototype-branches=branch-1 branch-2
```
