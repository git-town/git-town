# git town status reset

<a type="command-summary">

```command-summary
git town status reset [-h | --help] [-v | --verbose]
```

</a>

The _status reset_ command deletes the persisted runstate. This is only needed
if the runstate is corrupted and causes Git Town to crash.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [status show](status-show.md) displays the runstate that this command would
  delete

<!-- keep-sorted end -->
