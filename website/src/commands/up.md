# git town up

<a type="gittown-command" />

```command-summary
git town up [(-d | --display-types) <type>] [-h | --help] [-m | --merge] [(-o | --order) <asc|desc>] [-v | --verbose]
```

The _up_ command moves one position up in the current stack by switching to a
child of the current branch. After successfully switching branches, it displays
the branch hierarchy to show your new position in the stack.

When the current branch has multiple children, an interactive dialog lets you
choose which child branch to switch to.

`git town up` is useful for navigating stacked changes without needing to
remember branch names or use the interactive [switch](switch.md) command.

## Examples

Consider this stack:

```
main
 \
* branch-1
   \
    branch-2
```

After running `git town up` on the `branch-1` branch, you end up with this
stack:

```
main
 \
  branch-1
   \
*   branch-2
```

## Options

#### `-d`<br>`--display-types`

This flag allows customizing whether Git Town also displays the branch type in
addition to the branch name when showing a list of branches. More info
[here](../preferences/display-types.md#cli-flags).

#### `-m`<br>`--merge`

The `--merge` aka `-m` flag has the same effect as the
[git checkout -m](https://git-scm.com/docs/git-checkout#Documentation/git-checkout.txt--m)
flag. It attempts to merge uncommitted changes in your workspace into the target
branch.

This is useful when you have uncommitted changes in your current branch and want
to move them up to a child branch.

#### `-o`<br>`--order`

The `--order` flag allows customizing the order in which branches get displayed.
More info [here](../preferences/order.md#cli-flag)

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [branch](branch.md) displays the branch hierarchy
- [down](down.md) moves one position down in the current stack
- [swap](swap.md) changes the stack by swapping the position of current branch
  with its parent
- [switch](switch.md) interactively switch between branches

<!-- keep-sorted end -->
