# git rename-branch [old name] &lt;new name&gt;

The rename-branch command changes the name of the current branch in the local
and origin repository. It aborts if the new branch name already exists or the
tracking branch is out of sync.

### Variations

Provide the additional `old_name` argument to rename the branch with the given
name instead of the currently checked out branch. Renaming perennial branches
requires confirmation with the `-f` option.
