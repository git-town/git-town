# git town status

> _git town status [--pending]_

The _status_ command indicates whether Git Town has encountered a merge conflict
and which commands you can run to continue, skip, or undo it.

### --pending / -p

The `--pending` aka `-p` argument causes this command to output only the name of
the pending Git Town command if one exists. This allows displaying a reminder to
run `git town continue` into your shell prompt when you encountered a merge
conflict earlier. See [Integration](../integration.md#shell-prompt) on how to
set this up.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

### subcommands

Running without a subcommand shows the current Git Town configuration.

- The `get-parent` subcommand prints the parent branch of the current or given
