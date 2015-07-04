#### NAME

git-hack - cut a new feature branch off the main branch


#### SYNOPSIS

```
git hack <branchname> [parent branch name]
git hack (--abort | --continue)
```


#### DESCRIPTION

Syncs the main branch if there is a remote repository.
Creates a new feature branch with the given name.
Brings over all uncommitted changes.


#### OPTIONS

```
<branchname>
    The name of the branch to create.

[parent branch name]
    If provided, cuts the new branch off the given existing feature branch.
    Providing '.' here uses the current branch as the parent branch.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
