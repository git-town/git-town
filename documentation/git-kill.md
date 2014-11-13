#### NAME

git-kill - completely remove an obsolete feature branch


#### SYNOPSIS

```
git kill [<branchname>]
git kill --undo
```


#### DESCRIPTION

Deletes the current/given branch from the local and remote repository.


#### OPTIONS

```
<branchname>
    The branch to ship.
    If not provided, uses the current branch.

--undo
    Restore the last deleted Git branch.
```
