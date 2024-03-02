# git ship [branch name] [-m message]

The _ship_ command ("let's ship this feature") merges a completed feature branch
into the main branch and removes the feature branch. After the merge it pushes
the main branch to share the new commit on it with the rest of the world.

Git ship opens the default editor with a prepopulated commit message that you
can modify. You can submit an empty commit message to abort the shipping
process.

This command ships only direct children of the main branch. To ship a child
branch, you need to first ship or [kill](kill.md) all its ancestor branches.

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

If [sync-before-ship](../preferences/sync-before-ship.md) is enabled, Git Town
syncs the current branch before executing the ship. This allows you to resolve
merge conflicts on the feature branch instead of on the main branch. This helps
keep the main branch green, but can delay shipping.
