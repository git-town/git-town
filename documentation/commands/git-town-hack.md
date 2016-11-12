#### NAME

git-town-hack - create a new feature branch


#### SYNOPSIS

```
git town-hack <branch_name> [<parent_branch_name>]
git town-hack (--abort | --continue)
```


#### DESCRIPTION

Syncs the given parent branch (default: main branch),
forks a new feature branch with the given name off it,
pushes the new feature branch to the remote repository,
and brings over all uncommitted changes to the new feature branch.

Additionally, when there is a remote upstream,
the main branch is synced with its upstream counterpart.
This can be disabled by toggling the "hack-push-flag" configuration:

```
git town hack-push-flag false
```


#### OPTIONS

```
<branch_name>
    The name of the branch to create.

<parent_branch_name>
    If provided, cuts the new branch off the given existing feature branch.
    Providing '.' here uses the current branch as the parent branch.
    If omitted, uses the main branch as the parent.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
