# Display the branch stack in proposals

Git Town can embed a visual breadcrumb in proposals (pull requests, merge
requests), showing where the branch under review sits in your stack.

![example stack created by the Git Town GitHub action](https://raw.githubusercontent.com/git-town/action/main/docs/example-visualization.png)

These branch stacks get automatically updated when:

- proposing a branch
- shipping a branch
- prepending a branch
- detaching a branch
- merging branches
- swapping branches

You have two options to maintain such stacks. You only need to enable one of
them.

### Configure the Git Town executable to create and update branch stacks

The Git Town executable itself can maintain branch stacks in proposals. The
advantage of using Git Town for this is that this works with all forge types and
doesn't require changes to the CI system.

To make Git Town display branch stacks in proposals, configure these settings:

- [proposals-show-lineage](../preferences/proposals-show-lineage.md) enables or
  disables branch lineage in proposals

### Set up the Git Town GitHub action

If your entire team has standardized on using Git Town for branch management,
and you use GitHub, you can set up the
[Git Town GitHub action](https://github.com/marketplace/actions/git-town-github-action)
to automatically add branch stacks to pull requests.
