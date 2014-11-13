#### NAME

git-ship - ship a completed feature branch

#### SYNOPSIS

```
git ship [<branchname>] [-m <message>]
git ship --abort
```

#### DESCRIPTION

Squash merges the current branch, or `<branchname>` if given, into the main branch, leading to linear history on the main branch.

* sync the main branch
* pull remote updates for `<branchname>`
* merges the main branch into `<branchname>`
* squash-merges `<branchname>` into the main branch
* pushes the main branch to the remote repository
* deletes `<branchname>` from the local and remote repositories


#### OPTIONS

```
<branchname>
    The branch to ship.
    If not provided, uses the current branch.

-m <message>
    The commit message for the squash merge.
    If not provided, will be prompted for it.

--abort
    Cancel the operation and reset the workspace to a consistent state.
```
