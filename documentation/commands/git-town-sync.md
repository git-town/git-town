#### NAME

git-town-sync - updates the current branch with all relevant changes


#### SYNOPSIS

```
git town-sync [--all]
git town-sync (--abort | --continue | --skip)
```

#### DESCRIPTION

Synchronizes the current branch with the rest of the world.

When run on a feature branch
* syncs all ancestor branches
* pulls updates for the current branch
* merges the parent branch into the current branch
* pushes the current branch

When run on the main branch or a perennial branch
* pulls and pushes updates for the current branch
* pushes tags

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.

#### OPTIONS

```
--all
    Syncs all local branches

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.

--skip
    Continue the operation by skipping the sync of the current branch.
```
