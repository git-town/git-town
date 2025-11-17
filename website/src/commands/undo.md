# git town undo

<a type="command-summary">

```command-summary
git town undo [-h | --help] [-v | --verbose]
```

</a>

The _undo_ command reverts the last fully executed Git Town command. It performs
the opposite activities that the last command did and leaves your repository in
the state it was before you ran the problematic command.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [continue](continue.md) continues the currently suspended Git Town command
  after you have resolved the conflicting changes
- [skip](skip.md) ignores all remaining merge conflicts on the current branch
  and then continues the currently suspended Git Town command
