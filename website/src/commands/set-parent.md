# git town set-parent

> _git town set-parent_

The _set-parent_ command changes the parent branch for the current branch. You
select the new parent through a visual dialog. Updates associated proposals and
removes commits from former parent branches.

Since set-parent changes commits, you are strongly advised to sync your branches
before running this command.

After set-parent runs, the affected branches no longer contain changes made by
their old parents. However, they don't see the changes made by their new parent
branches yet. Please run [git town sync](sync.md) to pull in changes from the
new parents.

To demonstrate how `git town set-parent works`, let's say we have this branch
hierarchy:

```
main
 \
  feature-1
   \
*   feature-2
```

`feature-1` is a child branch of `main`, and `feature-2` is a child branch of
`feature-1`. Assuming we are on `feature-2`, we can make `feature-2` a child of
`main` by running `git town set-parent` and selecting `main` in the dialog. We
end up with this branch hierarchy:

```
main
 \
  feature-1
 \
* feature-2
```

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
