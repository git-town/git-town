# git town down

```command-summary
git town down [-m | --merge] [-v | --verbose]
```

The _down_ command moves one position down in the current stack by switching to
the parent of the current branch. After successfully switching branches, it
displays the branch hierarchy to show your new position in the stack.

`git town down` is useful for navigating stacked changes without needing to
remember branch names or use the interactive [switch](switch.md) command.

## Examples

Consider this stack:

```
main
 \
  branch-1
   \
*   branch-2
```

After running `git town down` on the `branch-2` branch, you end down with this
stack:

```
main
 \
* branch-2
   \
    branch-1
```

## Options

#### `-m`<br>`--merge`

The `--merge` aka `-m` flag has the same effect as the
[git checkout -m](https://git-scm.com/docs/git-checkout#Documentation/git-checkout.txt--m)
flag. It attempts to merge uncommitted changes in your workspace into the target
branch.

This is useful when you have uncommitted changes in your current branch and want
to move them down to the parent branch.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [branch](branch.md) displays the branch hierarchy
- [switch](switch.md) interactively switch between branches
- [swap](swap.md) changes the stack by swapping the position of current branch
  with its parent
- [up](up.md) moves one position up in the current stack
