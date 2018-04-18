<a textrun="command-heading">
# Append command
</a>

<a textrun="command-summary">
Creates a new feature branch as a direct child of the current branch.
</a>

<a textrun="command-description">
Syncs the current branch,
forks a new feature branch with the given name off the current branch,
makes the new branch a child of the current branch,
pushes the new feature branch to the remote repository
if and only if [new-branch-push-flag](./new-branch-push-flag.md) is true,
and brings over all uncommitted changes to the new feature branch.
</a>

#### Usage

<a textrun="command-cli">
```
git town append <branch>
```
</a>
