# git town set-parent

```command-summary
git town set-parent [<new-parent>] [-v | --verbose]
```

The _set-parent_ command changes the parent branch for the current branch.
Consider this stack:

```
main
 \
  feature-1
   \
*   feature-B
 \
  feature-A
```

After running `git town set-parent` and selecting `feature-A` in the dialog, we
end up with this stack:

```
main
 \
  feature-1
 \
  feature-A
   \
*   feature-B
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

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
