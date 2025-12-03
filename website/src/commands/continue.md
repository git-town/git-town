# git town continue

<a type="git-town-command" />

```command-summary
git town continue [-h | --help] [-v | --verbose]
```

When a Git Town command encounters a problem that it cannot resolve, for example
a merge conflict, it stops to give the user an opportunity to resolve the issue.
Once you have resolved the issue, run the _continue_ command to tell Git Town to
continue executing the failed command. Git Town will retry the failed operation
and execute all remaining operations of the original command.

## Options

#### `-h`<br>`--help`

Display help for this command.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [skip](skip.md) ignores all remaining merge conflicts on the current branch
  and then continues the currently suspended Git Town command
- [undo](undo.md) aborts the currently suspended Git Town command and returns
  the repository to the state it was in before you ran that command

<!-- keep-sorted end -->
