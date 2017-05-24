#### NAME

append - create a new feature branch as a direct child of the current branch


#### SYNOPSIS

```
git town append <branch_name>
git town append (--abort | --continue)
```


#### DESCRIPTION

Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the remote repository if and only if "hack-push-flag" is true,
and brings over all uncommitted changes to the new feature branch.


#### OPTIONS

```
<branch_name>
    The name of the branch to create.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
