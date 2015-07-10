#### NAME

git-ship - deliver a completed feature branch


#### SYNOPSIS

```
git ship [<branchname>] [<commit-options>]
git ship (--abort | --continue)
```


#### DESCRIPTION

Squash-merges the current branch, or `<branchname>` if given,
into the main branch, leading to linear history on the main branch.

* syncs the main branch
* pulls remote updates for `<branchname>`
* merges the main branch into `<branchname>`
* squash-merges `<branchname>` into the main branch
* pushes the main branch to the remote repository
* deletes `<branchname>` from the local and remote repositories

Only shipping of direct children of the main branch is allowed.
To ship a nested child branch, all ancestor branches have to be shipped or killed.


#### OPTIONS

```
<branchname>
    The branch to ship.
    If not provided, uses the current branch.

<commit-options>
    Options to pass to 'git commit' when commiting the squash-merge.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
