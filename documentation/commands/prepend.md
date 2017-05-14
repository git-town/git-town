#### NAME

prepend - create a new feature branch between the current branch and its parent


#### SYNOPSIS

```
git town prepend <branch_name>
git town prepend (--abort | --continue)
```


#### DESCRIPTION

Syncs the parent branch (prompts if unknown),
forks a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
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

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```


#### SEE ALSO
* [git append](append.md) to create a new feature branch as a child of the current branch
* [git hack](hack.md) to create a new top-level feature branch
