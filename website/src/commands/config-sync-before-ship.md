# git town config sync-before-ship [(yes|no)]

The _sync-before-ship_ configuration command displays or updates the
sync-before-ship configuration setting. If set to `yes`, [ship](ship.md)
executes [git sync](sync.md) before shipping a branch. This allows you to deal
with breakage from resolving merge conflicts on the feature branch instead of
the main branch. The downside is that this will trigger another CI run.

### Arguments

By default, each Git repository has its own setting. The `--global` flag
displays or sets "sync-before-ship" for all Git repos on your machine.
