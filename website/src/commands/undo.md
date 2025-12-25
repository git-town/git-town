# git town undo

<a type="git-town-command" />

```command-summary
git town undo [-h | --help] [-v | --verbose]
```

The _undo_ command reverts the last fully executed Git Town command. It performs
the opposite activities that the last command did and leaves your repository in
the state it was before you ran the problematic command.

## Options

#### `-h`<br>`--help`

Display help for this command.

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
