# git town undo

```command-summary
git town undo [-v | --verbose]
```

The _undo_ command reverts the last fully executed Git Town command. It performs
the opposite activities that the last command did and leaves your repository in
the state it was before you ran the problematic command.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
