# Run pre-push hook

This setting determines whether Git Town allows or prevents Git hooks while
pushing branches. Hooks are enabled by default. If your Git hooks are slow, you
can disable them to speed up branch syncing.

When disabled, Git Town pushes using the
[--no-verify](https://git-scm.com/docs/git-push) option. This omits the
[pre-push](https://git-scm.com/docs/githooks#_pre_push) hook.

## config file

To configure running the push hook in the
[configuration file](../configuration-file.md):

```toml
[sync]
push-hook = true
```

## Git metadata

To configure running the push hook manually in Git, run this command:

```wrap
git config [--global] git-town.push-hook <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the push hook by setting the `GIT_TOWN_PUSH_HOOK` environment
variable.
