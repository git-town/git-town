# git town prepend

> _git town prepend [--prototype] &lt;branch-name&gt;_

The _prepend_ command creates a new feature branch as the parent of the current
branch. It does that by inserting the new feature branch between the current
feature branch and it's existing parent.

When running without uncommitted changes in your workspace, it also
[syncs](sync.md) the current feature branch to ensure commits into the new
branch are on top of the current state of the repository. If the workspace
contains uncommitted changes, `git prepend` does not perform this sync to let
you commit your open changes first and then sync manually.

Consider this branch setup:

```
main
 \
* feature-2
```

We are on the `feature-2` branch. After running `git prepend feature-1`, our
repository has this branch setup:

```
main
 \
* feature-1
   \
    feature-2
```

### --prototype / -p

Adding the `--prototype` aka `-p` switch creates a
[prototype branch](../branch-types.md#prototype-branches)).

### Configuration

If [push-new-branches](../preferences/push-new-branches.md) is set, `git hack`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git hack` run fast. The first run of `git sync`
will create the remote tracking branch.

If the configuration setting
[create-prototype-branches](../preferences/create-prototype-branches.md) is set,
`git prepend` always creates a
[prototype branch](../branch-types.md#prototype-branches).
