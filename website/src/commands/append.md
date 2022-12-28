# git append &lt;branch&gt;

The _append_ command creates a new feature branch with the given name as a
direct child of the current branch and brings over all uncommitted changes to
the new branch. Before it does that, it [syncs](sync.md) the current branch to
ensure commits into the new branch are on top of the current state of the
repository.

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

### Variations

If [new-branch-push-flag](config-new-branch-push-flag.md) is set, `git append`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git append` run fast. The first run of `git sync`
will create the remote tracking branch.
