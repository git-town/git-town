# git town append

> _git town append [--prototype] &lt;branch-name&gt;_

The _append_ command creates a new feature branch with the given name as a
direct child of the current branch and brings over all uncommitted changes to
the new branch.

When running without uncommitted changes in your workspace, it also
[syncs](sync.md) the current branch to ensure your work in the new branch
happens on top of the current state of the repository. If the workspace contains
uncommitted changes, `git town append` does not perform this sync to let you
commit your open changes first and then sync manually.

### Positional argument

When given a non-existing branch name, `git town append` creates a new feature
branch with the main branch as its parent.

Consider this branch setup:

```
main
 \
* feature-1
```

We are on the `feature-1` branch. After running `git town append feature-2`, our
repository will have this branch setup:

```
main
 \
  feature-1
   \
*   feature-2
```

### --detached / -d

The `--detached` aka `-d` flag does not pull updates from the main or perennial
branch. This allows you to build out your branch stack and decide when to pull
in changes from other developers.

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

If [push-new-branches](../preferences/push-new-branches.md) is set,
`git town append` also creates the tracking branch for the new feature branch.
This behavior is disabled by default to make `git town append` run fast and save
CI runs. The first run of `git town sync` will create the remote tracking
branch.

If the configuration setting
[create-prototype-branches](../preferences/create-prototype-branches.md) is set,
`git town append` always creates a
[prototype branch](../branch-types.md#prototype-branches).
