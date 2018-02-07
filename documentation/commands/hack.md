#### NAME

hack - create a new feature branch off the main development branch

#### SYNOPSIS

```
git town hack <branch_name>
git town hack (--abort | --continue)
```

#### DESCRIPTION

Syncs the main branch and forks a new feature branch with the given name off it.

If (and only if) [new-branch-push-flag](./new-branch-push-flag.md) is true,
pushes the new feature branch to the remote repository.

Finally, brings over all uncommitted changes to the new feature branch.

#### OPTIONS

```
<branch_name>
    The name of the branch to create.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
