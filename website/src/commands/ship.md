# git ship [branch name] [-m message]

_Notice: Most people don't need to use the _ship_ command. The recommended way
to merge your feature branches is to use the web UI or merge queue of your code
hosting service, as you would normally do. `git ship` is for edge cases like
developing in [offline mode](../preferences/offline.md)._

The _ship_ command ("let's ship this feature") squash-merges a completed feature
branch into the main branch and removes the feature branch. After the merge it
pushes the main branch to share the new commit on it with the rest of the world.

Git ship opens the default editor with a prepopulated commit message that you
can modify. You can submit an empty commit message to abort the shipping
process.

This command ships only direct children of the main branch. To ship a child
branch, you need to first ship or [kill](kill.md) all its ancestor branches. If
you really need to ship into a non-perennial branch, you can override the
protection against that with the `--to-parent` option.

The branch to ship must be in sync. If it isn't in sync, `git ship` will exit
with an error. When that happens, run [git sync](sync.md) to get the branch in
sync, re-test and re-review the updated branch, and then run `git ship` again.

### Arguments

Similar to `git commit`, the `-m` parameter allows specifying the commit message
via the CLI.

### Configuration

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
