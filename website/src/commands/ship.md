# git town ship

> _git town ship [--to-parent] [--message &lt;text&gt;] [branch-name]_

_Notice: Most people don't need to use the _ship_ command. The recommended way
to merge your feature branches is to use the web UI or merge queue of your code
hosting service, as you would normally do. `git ship` is for edge cases like
developing in [offline mode](../preferences/offline.md) or when shipping
[stacked changes](../stacked-changes.md)._

The _ship_ command ("let's ship this feature") merges a completed feature branch
into the main branch and removes the feature branch.

The branch to ship must be in sync. If it isn't in sync, `git ship` will exit
with an error. When that happens, run [git sync](sync.md) to get the branch in
sync, re-test and re-review the updated branch, and then run `git ship` again.

### Positional argument

When called without a positional argument, the _ship_ command ships the current
branch.

When called with a positional argument, it ships the branch with the given name.

### --dry-run

The `--dry-run` flag allows to test-drive this command. It prints the Git
commands that would be run but doesn't execute them.

### --message / -m

Similar to `git commit`, the `--message <message>` aka `-m` parameter allows
specifying the commit message via the CLI.

### --to-parent / -p

The _ship_ command ships only direct children of the main branch. To ship a
child branch, you need to first ship or [kill](kill.md) all its ancestor
branches. If you really want to ship into a non-perennial branch, you can
override the protection against that with the `--to-parent` aka `-p` option.

### --verbose / -v

The `--verbose` aka `-v` flag prints all Git commands run under the hood to
determine the repository state.

### Configuration

The configured [ship-strategy](../preferences/ship-strategy.md) determines how
the _ship_ command merges branches. When shipping
[stacked changes](../stacked-changes.md), use the
[fast-forward ship strategy](../preferences/ship-strategy.md#fast-forward) to
avoid empty merge conflicts.

If you have configured the API tokens for
[GitHub](../preferences/github-token.md),
[GitLab](../preferences/gitlab-token.md), or
[Gitea](../preferences/gitea-token.md) and the branch to be shipped has an open
proposal, this command merges the proposal for the current branch on your origin
server rather than on the local Git workspace.

If your origin server deletes shipped branches, for example
[GitHub's feature to automatically delete head branches](https://help.github.com/en/github/administering-a-repository/managing-the-automatic-deletion-of-branches),
you can
[disable deleting remote branches](../preferences/ship-delete-tracking-branch.md).
