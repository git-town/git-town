# git town config ship-delete-tracking-branch [(yes|no)]

The _ship-delete-tracking-branch_ configuration command displays or updates the
ship-delete-tracking-branch configuration setting. If set to `yes`,
[ship](ship.md) deletes the tracking branch of shipped branches. Disable this if
your code hosting service deletes shipped branches on its end.

### Arguments

By default, each Git repository has its own setting. The `--global` flag
displays or sets "ship-delete-tracking-branch" for all Git repos on your
machine.
