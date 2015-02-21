#### NAME

git-extract - copy selected commits from the current branch into their own branch


#### SYNOPSIS

```
Usage:
  git extract <branchname> [<commit>...]
  git extract --abort | --continue
```


#### DESCRIPTION

If no commits are provided, prompts the user to select from a list of commits unique to the current branch.

* sync the main branch
* create a feature branch off it
* cherry pick commits


#### OPTIONS

```
<branchname>
    The name of the branch to create.

<commit>
    SHA to be cherry-picked into the new branch.

--abort
    Cancel the operation and reset the workspace to a consistent state.

--continue
    Continue the operation after resolving conflicts.
```
