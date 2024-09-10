# git town kill

> _git town kill [branch name]_

The _kill_ command deletes non-perennial branches from the local and remote
repository. It does not delete perennial branches.

When killing the currently checked out branch, you end up on the previously
checked out branch. If that branch also doesn't exist, you end up on the main
development branch.

### Positional arguments

When called without arguments, the _kill_ command deletes the feature branch you
are on including all uncommitted changes.

When called with a branch name, it kills the given branch.
