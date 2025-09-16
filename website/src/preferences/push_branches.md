# push-branches

This setting determines whether Git Town pushes local branches and commits to
the [development remote](dev-remote.md).

## config file

```toml
[sync]
push-branches = true
```

## Git metadata

To configure whether branches get pushed manually in Git, run this command:

```wrap
git config [--global] git-town.push-branches <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure whether branches get pushed by setting the
`GIT_TOWN_PUSH_BRANCHES` environment variable.
