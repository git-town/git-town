# git town set-parent

<a type="git-town-command" />

```command-summary
git town set-parent [<branch>] [--auto-resolve] [-h | --help] [--none] [-v | --verbose]
```

The _set-parent_ command moves a branch and all its children below another
branch. Consider this stack:

```
main
 \
  feature-1
   \
*   feature-2
     \
      feature-3

 \
  feature-A
```

Running `git town set-parent feature-A` creates this stack:

```
main
 \
  feature-1

 \
  feature-A
   \
*   feature-2
     \
      feature-3
```

You can also use `set-parent` to make a child branch a sibling branch. Consider
this stack:

```
main
 \
  feature-1
   \
    feature-2
     \
*     feature-3
       \
        feature-4
```

Running `git town set-parent feature-1` creates this stack:

```
main
 \
  feature-1
   \
    feature-2
   \
*   feature-3
     \
      feature-4
```

Since set-parent changes commits, your branches must be in sync when running
this command. Run [git town sync](sync.md) before running `git town set-parent`.

After set-parent runs, the affected branches no longer contain changes made by
their old parents. However, they don't see the changes made by their new parent
branches yet. Please run [git town sync](sync.md) to pull in changes from the
new parents.

## Positional argument

You can provide the name of the new parent for the current branch as an argument
to this command. When called without arguments, queries the user for the new
parent.

## Options

#### `--auto-resolve`

Disables automatic resolution of
[phantom merge conflicts](../stacked-changes.md#avoid-phantom-conflicts).

#### `-h`<br>`--help`

Display help for this command.

#### `--none`

The `--none` option assigns no parent (removes the assigned parent), making the
branch a perennial branch.

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

## See also

<!-- keep-sorted start -->

- [detach](detach.md) extract the current branch from a stack, leaving its
  children in the stack.
- [swap](swap.md) move the current branch up in the stack

<!-- keep-sorted end -->
