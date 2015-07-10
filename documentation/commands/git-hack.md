#### NAME

git-hack - create a new feature branch


#### SYNOPSIS

```
git hack <branchname> [<parent branch name>]
git hack (--abort | --continue)
```


#### DESCRIPTION

Syncs the given parent branch,
forks a new feature branch off it,
and brings over all uncommitted changes to the new feature branch.


#### OPTIONS

```
<branchname>
    The name of the branch to create.

[parent branch name]
    If provided, cuts the new branch off the given existing feature branch.
    Providing '.' here uses the current branch as the parent branch.
    If omitted, uses the main branch as the parent.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
