# git town offline

> _git town offline <status>_

The _offline_ configuration command displays or changes Git Town's offline mode.
Git Town skips all network operations in offline mode.

### Positional arguments

When called without an argument, the _offline_ command displays the current
offline status.

When called with `yes`, `1`, `on`, or `true`, this command enables offline mode.
When called with `no`, `0`, `off`, or `false`, it disables offline mode.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
