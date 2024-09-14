# git town kill

> _git town kill [branch-name]_

The _kill_ command deletes non-perennial branches from the local and remote
repository and updates proposals of child branches of the killed branch. It does
not remove perennial branches.

### Positional arguments

When called without arguments, the _kill_ command deletes the feature branch you
are on, including all uncommitted changes.

When called with a branch name, it kills the given branch.

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
