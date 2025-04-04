# Create prototype branches

_This setting is deprecated. Please replace it with the
[new-branch-type](new-branch-type.md) setting._

This setting determines whether Git Town creates new branches as
[prototype branches](../branch-types.md#prototype-branches).

## config file

To configure the creation of prototype branches in the
[configuration file](../configuration-file.md):

```toml
create-prototype-branches = true
```

## Git metadata

To configure the creation of prototype branches manually in Git, run this
command:

```wrap
git config [--global] git-town.create-prototype-branches <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
