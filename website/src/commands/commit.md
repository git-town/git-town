# git town commit

<a type="git-town-command" />

```command-summary
git town commit [-d | --down] [--dry-run] [-h | --help] [(-m | --message) <text>] [-v | --verbose]
```

The _commit_ command takes the currently staged changes and commits them into a
different branch in your stack, then synchronizes the result back into your
current branch.

This is useful when working with [stacked branches](../stacked-changes.md). A
common scenario is that you're implementing a feature and realize that part of
the work really belongs in a refactor. You want that refactor reviewed and
shipped independently, but the feature depends on it, so the refactor must live
in an ancestor branch.

Your desired branch stack might look like this:

```
main
 \
  refactor
   \
    feature
```

Manually switching back and forth between `refactor` and `feature` to move
changes around is slow and error-prone.

With `git town commit`, you can stay on the `feature` branch, do the refactoring
there, and then commit those changes directly into the `refactor` branch. Git
Town will automatically sync the committed changes back into `feature`, letting
you continue where you left off.

## Options

#### `-d`<br>`--down`

Commit the staged changes into the parent branch of the current branch.

#### `--dry-run`

Print the Git commands that would be executed without actually running them.

#### `-h`<br>`--help`

Display help for this command.

#### `-m <text>`<br>`--message <text>`

Set the commit message from the command line, equivalent to `git commit -m`.

#### `-v`<br>`--verbose`

Prints all Git commands executed under the hood, used to determine repository
state.

## See also

- [git town prepend --commit](prepend.md#-c--commit)
- [git town prepend --beam](prepend.md#-b--beam)
- [git town append --commit](append.md#-c--commit)
- [git town append --beam](append.md#-b--beam)
- [git town hack --commit](hack.md#-c--commit)
- [git town hack --beam](hack.md#-b--beam)
