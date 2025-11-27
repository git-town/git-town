# git town status

<a type="gittown-command" />

```command-summary
git town status [-h | --help] [-p | --pending] [-v | --verbose]
```

The _status_ command indicates whether Git Town has encountered a merge conflict
and which commands you can run to continue, skip, or undo it.

## Subcommands

The [reset](status-reset.md) subcommand deletes the persisted runstate. This is
only needed if the runstate is corrupted and causes Git Town to crash.

The [show](status-show.md) subcommand displays detailed information about the
persisted runstate.

## Options

#### `-h`<br>`--help`

Display help for this command.

#### `-p`<br>`--pending`

The `--pending` aka `-p` argument causes this command to output only the name of
the pending Git Town command if one exists. This allows displaying a reminder to
run `git town continue` into your shell prompt when you encountered a merge
conflict earlier. See [here](../how-to/shell-prompt.md) on how to set this up

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
