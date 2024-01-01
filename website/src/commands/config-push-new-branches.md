# git town config push-new-branches [(true|false)]

The _push-new-branches_ configuration command displays or updates the
push-new-branches configuration setting. If set to `yes`, [hack](hack.md),
[append](append.md), and [prepend](prepend.md) push newly created feature
branches to the `origin` remote. Defaults to `no`.

### Arguments

By default, each Git repository has its own setting. The `--global` flag
displays or sets the "push-new-branches" for all Git repos on your machine.
