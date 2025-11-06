# git town branch

```command-summary
git town branch [-d | --display-types] [-o | --order] [-v | --verbose]
```

The _branch_ command is Git Town's equivalent of the
[git branch](https://git-scm.com/docs/git-branch) command. It displays the local
branch hierarchy, and the types of all branches except for main and feature
branches.

## Options

#### `-d`<br>`--display-types`

This flag allows customizing whether Git Town also displays the branch type in
addition to the branch name when showing a list of branches. More info
[here](../preferences/display-types.md#cli-flags).

#### `-o`<br>`--order`

The `--order` flag allows customizing the order in which branches get displayed.
More info [here](../preferences/order.md#cli-flag)

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

- [switch](switch.md) displays the branch hierarchy and lets you switch to a new
  branch in it
- [walk](walk.md) executes a shell command or opens a shell in each of your
  local branches
