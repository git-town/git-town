# git append &lt;branch&gt;

The _append_ command creates a new feature branch with the given name as a
direct child of the current branch and brings over all uncommitted changes to
the new branch.

When running without uncommitted changes in your workspace, it also
[syncs](sync.md) the current branch to ensure your work in the new branch
happens on top of the current state of the repository. If the workspace contains
uncommitted changes, `git append` does not perform this sync to let you commit
your open changes first and then sync manually.

### Example

Consider this branch setup:

```
main
 \
  feature-1
```

We are on the `feature-1` branch. After running `git append feature-2`, our
repository will have this branch setup:

```
main
 \
  feature-1
   \
    feature-2
```

### Configuration

If [push-new-branches](../preferences/push-new-branches.md) is set, `git append`
also creates the tracking branch for the new feature branch. This behavior is
disabled by default to make `git append` run fast. The first run of `git sync`
will create the remote tracking branch.
