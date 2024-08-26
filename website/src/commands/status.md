# git town status

The _status_ command indicates whether Git Town has encountered a merge conflict
and which commands you can run to continue, skip, or undo it.

### Arguments

The `--pending` argument causes this command to output only the name of the
pending Git Town command if one exists, and nothing otherwise. This allows
displaying a reminder to run `git town continue` into your shell prompt when you
encountered a merge conflict earlier. See
[Integration](../integration.md#shell-prompt) on how to set this up.
