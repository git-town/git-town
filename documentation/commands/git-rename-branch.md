#### NAME

git-rename-branch - rename a branch both locally and remotely


#### SYNOPSIS

```
git rename-branch <old_branch_name> <new_branch_name> [-f]
```


#### DESCRIPTION

Renames the given branch on both the local machine and the remote if one is configured.
Aborts if the new branch name already exists or the tracking branch is out of sync.
This command is intended for feature branches. Renaming perennial branches has to be confirmed with the `-f` option.

* Creates a branch with the new name
* Deletes the old branch

When there is a remote repository
* Syncs the repository

When there is a tracking branch
* Pushes the new branch to the remote repository
* Deletes the old branch from the remote repository

When run on a perennial branch
* Requires the use of the `-f` option
* Reconfigures git-town locally for the perennial branch


#### OPTIONS

```
<old_branch_name>
    The name of the branch to rename.

<new_branch_name>
    The new name of the branch.

-f
    Forces the renaming of a perennial branch
```
