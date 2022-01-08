# git town append &lt;branch&gt;

The append command creates a new feature branch with the given name as a direct
child of the current branch and brings over all uncommitted changes to the new
feature branch. Before it does that, it [syncs](sync.md) the current branch to
ensure the changes in the feature branch are on top of the latest code version.

### Customization

If [new-branch-push-flag](.new-branch-push-flag.md) is set, `git append` creates
a remote tracking branch for the new feature branch. This behavior is disabled
by default to make `git append` run fast. The first run of `git sync` will then
create the remote tracking branch.
