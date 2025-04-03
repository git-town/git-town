# git town set-parent

```command-summary
git town set-parent [-v | --verbose]
```

The _set-parent_ command changes the parent branch for the current branch. You
select the new parent through a visual dialog. Updates associated proposals and
removes commits from former parent branches.

Consider this branch stack:

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
end up with this branch stack:

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
this command.

After set-parent runs, the affected branches no longer contain changes made by
their old parents. However, they don't see the changes made by their new parent
branches yet. Please run [git town sync](sync.md) to pull in changes from the
new parents.

## Options

#### `-v`<br>`--verbose`

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
