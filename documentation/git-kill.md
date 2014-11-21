#### NAME

git-kill - remove an obsolete feature branch


#### SYNOPSIS

```
git kill [<branchname>]
git kill --undo
```


#### DESCRIPTION

Deletes the current branch, or `<branchname>` if given,
from the local and remote repositories.

Does not delete non-feature branches nor the main branch.



#### OPTIONS

```
<branchname>
    The branch to remove.
    If not provided, uses the current branch.

--undo
    undo the previous `git kill` operation
```
