# git town offline

<a type="git-town-command" />

```command-summary
git town offline [<status>] [-h | --help] [-v | --verbose]
```

The _offline_ configuration command displays or changes Git Town's offline mode.
Git Town skips all network operations in offline mode.

## Positional arguments

When called without an argument, the _offline_ command displays the current
offline status.

When called with `yes`, `1`, `on`, or `true`, this command enables offline mode.
When called with `no`, `0`, `off`, or `false`, it disables offline mode.

## Options

#### `-h`<br>`--help`

Display help for this command.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
