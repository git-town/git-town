# git town set-parent

> _git town set-parent_

The _set-parent_ command changes the parent branch for the current branch and
updates associated proposals. It prompts the user for the new parent branch.

This command does not update commits, i.e. the new child branches don't see the
changes made by their new parent branches. To update the commits, run
[git town sync](sync.md).

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
