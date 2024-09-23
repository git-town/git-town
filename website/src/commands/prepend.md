# git town prepend

> _git town prepend [--prototype] &lt;branch-name&gt;_

The _prepend_ command creates a new feature branch as the parent of the current
branch. It does that by inserting the new feature branch between the current
feature branch and it's existing parent.

If your Git workspace is clean (no uncommitted changes), it also
[syncs](sync.md) the current feature branch to ensure you work on top of the
current state of the repository. If the workspace is not clean (contains
uncommitted changes), `git town prepend` does not perform this sync to let you commit
your open changes.

If the branch you call this command from has a proposal, this command updates
it. To do so, it pushes the new branch.

Consider this branch setup:

```
main
 \
* feature-2
```

We are on the `feature-2` branch. After running `git town prepend feature-1`, our
repository has this branch setup:

```
main
 \
* feature-1
   \
    feature-2
```

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --prototype / -p

Adding the `--prototype` aka `-p` switch creates a
[prototype branch](../branch-types.md#prototype-branches)).

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

### Configuration

If [push-new-branches](../preferences/push-new-branches.md) is set, `git town hack`
creates a remote tracking branch for the new feature branch. This behavior is
disabled by default to make `git town hack` run fast. The first run of `git town sync`
will create the remote tracking branch.

If the configuration setting
[create-prototype-branches](../preferences/create-prototype-branches.md) is set,

[prototype branch](../branch-types.md#prototype-branches).
