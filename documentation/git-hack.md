#### NAME

git-hack - cut a new feature branch off the main branch

#### SYNOPSIS

```
git hack <branchname>
git hack -abort
```

#### DESCRIPTION

Sync the main branch and create a new feature branch with the given name.
Brings over all uncommitted changes.

#### OPTIONS

```
<branchname>
    The name of the branch to create.

--abort
    Cancel the operation and reset the workspace to a consistent state.
```
