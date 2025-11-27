# git town swap

<a type="gittown-command" />

```command-summary
git town swap [--auto-resolve] [--dry-run] [-h | --help] [-v | --verbose]
```

The _swap_ command switches the position of the current branch with the branch
ahead of it in the current stack, i.e. moves the current branch one position
forward in the stack.

Consider this stack:

```
main
 \
  branch-1
   \
*   branch-2
     \
      branch-3
```

After running `git town swap` on the `branch-2` branch, you end up with this
stack:

```
main
 \
  branch-2
   \
*   branch-1
     \
      branch-3
```

Moving branches up and down the stack allows you to organize related branches
together, for example to review and ship them as a series, or to
[merge](merge.md) them.

Please ensure that all affected branches are in sync and don't contain merge
commits before running this command, by running [git town sync](sync.md) and
optionally [git town compress](compress.md) before. All affected branches must
be owned by you, i.e. you cannot swap
[contribution](../branch-types.md#contribution-branches),
[observed](../branch-types.md#observed-branches), or
[perennial](../branch-types.md#perennial-branches) branches.

## Options

#### `--auto-resolve`

Disables automatic resolution of
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts).

#### `--dry-run`

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

#### `-h`<br>`--help`

Display help for this command.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [detach](detach.md) extracts the current branch from a stack, leaving its
  children in the stack.
- [set-parent](set-parent.md) moves the current branch and its descendents under
  a different parent

<!-- keep-sorted end -->
