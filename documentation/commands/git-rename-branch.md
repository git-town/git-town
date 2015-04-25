#### NAME

git-rename-branch - rename a branch both locally and remotely


#### SYNOPSIS

```
git rename-branch <branchname> <newbranchname> [-f]
```


#### DESCRIPTION

On a non-feature branch, requires the use of the `-f` option
Syncs the repository if there is a remote repository
Creates a branch with the new name based on the old name
If there is a remote repository
* pushes the new branch
* deletes the old branch from the remote repository
Deletes the old branch locally
Reconfigures git-town locally if renaming a non-feature branch



#### OPTIONS

```
<branchname>
    The name of the branch to rename.

<newbranchname>
    The new name of the branch being renamed.

-f
    Forces the renaming of a non-feature branch 
```
