#### NAME

git-sync - synchronize the current branch

#### SYNOPSIS

```
git sync
git sync --abort
git sync --continue
```

#### DESCRIPTION

Synchronizes the current branch with the rest of the world.

When run on a feature branch
* syncs the main branch
* pulls updates for the current branch
* merges the main branch into the current branch
* pushes the current branch

When run on the main branch or a non-feature branch
* pulls and pushes updates for the current branch
* pushes tags

#### OPTIONS

```
--abort
  Cancel the operation and reset the workspace to a consistent state.

--continue
  Continue the operation after resolving conflicts.
```
