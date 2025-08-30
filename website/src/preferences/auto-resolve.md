# auto-resolve

Git Town automatically resolves
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts) by
default. If you run into a bug related to auto-resolution, please
[report it](https://github.com/git-town/git-town/issues) and then use this
configuration setting to disable auto-resolution as a temporary workaround until
the issue is addressed.

## CLI flags

In one-off situations you can call commands that update branches with the
`--no-auto-resolve` flag to disable automatic resolution of phantom merge
conflicts.

If you have automatic resolution disabled permanently via the config file or Git
metadata (see below), you can enable it with the `--auto-resolve` flag.

## config file

To configure automatic resolution of phantom merge conflicts in the
[configuration file](../configuration-file.md):

```toml
[sync]
auto-resolve = false
```

## Git metadata

To configure phantom merge resolution in Git, run this command:

```wrap
git config [--global] git-town.auto-resolve <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure auto-resolving of phantom merge conflicts by setting the
`GIT_TOWN_AUTO_RESOLVE` environment variable.
