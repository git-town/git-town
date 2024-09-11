# git town kill

> _git town kill [branch-name]_

The _kill_ command deletes non-perennial branches from the local and remote
repository. It does not delete perennial branches.

When killing the currently checked out branch, you end up on the previously
checked out branch. If that branch also doesn't exist, you end up on the main
development branch.

### Positional arguments

When called without arguments, the _kill_ command deletes the feature branch you
are on including all uncommitted changes.

When called with a branch name, it kills the given branch.

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
