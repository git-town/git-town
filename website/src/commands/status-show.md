# git town status show

<a type="command-summary">

```command-summary
git town status show [-h | --help] [-v | --verbose]
```

</a>

The _status show_ command displays Git Town's runstate, i.e. detailed
information about the currently suspended or previously executed Git Town
command, including its path on the filesystem.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [runlog](runlog.md) displays an overview of the most recently executed Git
  Town commands
- [status reset](status-reset.md) deletes the runstate. This can solve errors
  after upgrading Git Town.

<!-- keep-sorted end -->
