# Development remote

This setting lets you override the name of the Git remote used for development.
This is the remote that branches get pushed to, and into which branches get
shipped to. Usually that remote is called `origin`, which is also the default
value for this setting.

## config file

To configure the development remote in the
[configuration file](../configuration-file.md):

```toml
[hosting]
dev-remote = "<remote name>"
```

## Git metadata

To configure the development remote manually in Git, run this command:

```wrap
git config [--global] git-town.dev-remote <remote name>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the development remote by setting the `GIT_TOWN_DEV_REMOTE`
environment variable.
