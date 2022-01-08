# git rename-branch [old name] &lt;new name&gt;

The rename-branch command changes the name of the current branch in the local
and origin repository. It aborts if the new branch name already exists or the
tracking branch is out of sync.

### Customization

Provide the `old_name` argument to rename a branch that is not currently checked
out. Confirm renaming perennial branches with the `-f` option.
