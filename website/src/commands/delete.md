# git town delete

> _git town delete [branch-name]_

The _delete_ command deletes the given branch from the local and remote
repository and updates proposals of its child branches to the parent of the
deleted branch. It does not remove perennial branches.

Removes commits of deleted branches from their descendents, unless when using
the [merge sync strategy](../preferences/sync-feature-strategy.md#merge).

### Positional arguments

When called without arguments, the _delete_ command deletes the feature branch
you are on, including all uncommitted changes.

When called with a branch name, it deletes the given branch.

### --dry-run

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
