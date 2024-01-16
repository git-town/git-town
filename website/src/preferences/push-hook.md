# push-hook configuration setting

The "push-hook" setting determines whether Git Town permits or prevents Git
hooks while pushing branches. By default, hooks are enabled. If your Git hooks
are slow, you can disable them to speed up the process of syncing branches.

When disabled, Git Town pushes using the
[--no-verify](https://git-scm.com/docs/git-push) option, which omits the
[pre-push](https://git-scm.com/docs/githooks#_pre_push) hook.

Usage in the [configuration file](../configuration-file.md):

```toml
push-hook = false
```
