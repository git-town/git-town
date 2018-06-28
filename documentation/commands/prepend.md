#### NAME

prepend - create a new feature branch between the current branch and its parent

#### SYNOPSIS

```
git town prepend <branch_name>
```

#### DESCRIPTION

Syncs the parent branch,
forks a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the remote repository
if and only if [new-branch-push-flag](./new-branch-push-flag.md) is true,
and brings over all uncommitted changes to the new feature branch.

#### OPTIONS

```
<branch_name>
    The name of the branch to create.
```

#### SEE ALSO

* [git append](append.md) to create a new feature branch as a child of the current branch
* [git hack](hack.md) to create a new top-level feature branch
