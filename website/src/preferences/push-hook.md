# push-hook configuration setting

The "push-hook" setting determines whether Git Town allows or prevents Git hooks
while pushing branches. Hooks are enabled by default. If your Git hooks are
slow, you can disable them to speed up branch syncing.

When disabled, Git Town pushes using the
[--no-verify](https://git-scm.com/docs/git-push) option. This omits the
[pre-push](https://git-scm.com/docs/githooks#_pre_push) hook.

The best way to change this setting is via the
[setup assistant](../configuration.md).

## config file

To configure the push hook in the
[configuration file](../configuration-file.md):

```toml
sync.push-hook = false
```

or

```toml
[sync]
push-hook = false
```

## Git metadata

To configure the push hook manually in Git, run this command:

```bash
git config [--global] git-town.push-hook <true|false>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.
