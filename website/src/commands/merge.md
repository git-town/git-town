# git town merge

> _git town merge_

The _merge_ command merges the current branch with its parent branch.
Both branches must be
[feature branches](../branch-types.md#feature-branches).

When using the
[compress sync strategy](../preferences/sync-feature-strategy.md#compress), the
merged branch will contain two separate commits: one per merged branch. This
makes it easy to verify that both branches were merged as expected. To
consolidate these commits, run [git town sync](sync.md).

### --dry-run

Use the `--dry-run` flag to test-drive this command. It prints the Git commands
that would be run but doesn't execute them.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.
