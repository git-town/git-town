# git town runlog

```command-summary
git town runlog [-v | --verbose]
```

Git Town records the SHA of all local and remote branches before and after each
command runs into an immutable, append-only log file called the runlog.

The runlog provides an extra layer of safety, making it easier to manually roll
back changes if [git town undo](undo.md) doesn’t fully undo the changes the last
command made.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
