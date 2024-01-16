# push-hook configuration setting

The "push-hook" setting determines whether Git Town permits or prevents Git
hooks while pushing branches. Hooks are enabled by default. If your Git hooks
are slow, you can disable them to speed up branch syncing.

When disabled, Git Town pushes using the
[--no-verify](https://git-scm.com/docs/git-push) option, which omits the
[pre-push](https://git-scm.com/docs/githooks#_pre_push) hook.

Usage in the [configuration file](../configuration-file.md):

```toml
push-hook = false
```
