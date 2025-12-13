# git town skip

<a type="git-town-command" />

```command-summary
git town skip [-h | --help] [--park] [-v | --verbose]
```

The _skip_ command allows to skip a Git branch with merge conflicts when syncing
all feature branches.

## Options

#### `-h`<br>`--help`

Display help for this command.

#### `--park`

When set, also [parks](park.md) the current branch.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [continue](continue.md) tries to continue the currently suspended Git Town
  command by re-running the Git command that failed.
- [undo](undo.md) aborts the currently suspended Git Town command and undoes all
  the changes it did so far, leaving your repository in the same state it was in
  before you started the failing Git Town command

<!-- keep-sorted end -->
