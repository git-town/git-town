# git town rename-branch

> _git town rename-branch [--force] [old-name] &lt;new-name&gt;_

The _rename-branch_ command changes the name of the current branch in the local
and origin repository. It aborts if the new branch name already exists or the
tracking branch is out of sync.

### Positional arguments

When called with only one argument, the _rename-branch_ command renames the
current branch to the given name.

When called with two arguments, it renames the branch with the given name to the
given name.

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --force / -f

Renaming perennial branches requires confirmation with the `--force` aka `-f`
flag.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
