# git town runlog

```command-summary
git town runlog [-v | --verbose]
```

Git Town logs the repository state before and after each Git Town command
executes. This is an additional safety net, which allows you to manually undo a
Git Town command in case [git town undo](undo.md) doesn't work fully.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
