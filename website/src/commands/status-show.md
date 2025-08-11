# git town status show

```command-summary
git town status show [-v | --verbose]
```

The _status show_ command displays Git Town's runstate, i.e. detailed
information about the currently suspended or previously executed Git Town
command, including its path on the filesystem.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [runlog](runlog.md) displays an overview of the most recently executed Git
  Town commands
- [status reset](status-reset.md) deletes the runstate. This can solve errors
  after upgrading Git Town.
