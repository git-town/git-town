# git town merge

> _git town merge_

The _merge_ command merges the current branch with its parent branch and updates
proposals. Both branches must be
[feature branches](../branch-types.md#feature-branches).

When using the
[compress sync strategy](../preferences/sync-feature-strategy.md#compress), the
merged branch will contain two commits: one commit per merged branch. This
allows you to verify that the branches were correctly merged. Run
[git town sync](sync.md) to compress these two commits.

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
