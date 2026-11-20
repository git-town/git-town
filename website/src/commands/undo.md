# git town undo

<a type="git-town-command" />

```command-summary
git town undo [-f | --force] [-h | --help] [--(non)-interactive] [-v | --verbose]
```

The _undo_ command reverts the last fully executed Git Town command.
It performs the opposite activities that the last command did
and leaves your repository in the state it was
before you ran the problematic command.

## Options

#### `-f`<br>`--force`

If a branch that the last Git Town command created received new commits since
then, undoing that command would delete these commits. In that case, the _undo_
command refuses to run. Use the `--force` aka `-f` flag to undo anyway.

#### `-h`<br>`--help`

Display help for this command.

#### `--interactive`<br>`--non-interactive`

Enables or disables [interactive](../preferences/interactive.md) dialogs.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [continue](continue.md) continues the currently suspended Git Town command
  after you have resolved the conflicting changes
- [skip](skip.md) ignores all remaining merge conflicts on the current branch
  and then continues the currently suspended Git Town command

<!-- keep-sorted end -->
