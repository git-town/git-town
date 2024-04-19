# git kill [branch]

The _kill_ command deletes the feature branch you are on including all
uncommitted changes from the local and remote repository. It does not delete
perennial branches.

When killing the currently checked out branch, you end up on the previously
checked out branch. If that branch also doesn't exist, you end up on the main
development branch.

### Arguments

If you provide an argument, `git kill` removes the branch with the given name
instead of the current branch.
