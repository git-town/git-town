# Display a breadcrumb in proposals

Git Town can render a visual breadcrumb in proposals that shows where the
current branch sits within its stack. This makes stacked changes explicit and
easier to review.

![example stack created by the Git Town GitHub action](https://raw.githubusercontent.com/git-town/action/main/docs/example-visualization.png)

These breadcrumbs are kept up to date automatically when you:

- [propose](../commands/propose.md) a branch
- [sync](../commands/sync.md) branches
- [ship](../commands/ship.md) a branch
- [delete](../commands/delete.md) a branch
- [prepend](../commands/prepend.md) a branch
- [detach](../commands/detach.md) a branch
- [change the parent](../commands/set-parent.md) of a branch
- [merge](../commands/merge.md) branches
- [swap](../commands/swap.md) branches

There are two ways to maintain breadcrumbs. You only need to enable one of them.

### Use the Git Town executable

The Git Town CLI can create and update breadcrumbs. This approach works across
all supported forges and doesn't require any CI or workflow changes. Breadcrumbs
only get updated if branch changes happen through the Git Town CLI.

To enable this behavior, set
[proposal-breadcrumb](../preferences/proposal-breadcrumb.md) to one of the
following values:

- `branches` to display breadcrumbs for all branches
- `stacks` to display breadcrumbs only for stacks that contain more than one
  branch
- `none` to not display breadcrumbs

### Use the GitHub action

If your team standardizes on Git Town and uses GitHub, you can set up the
[Git Town GitHub action](https://github.com/marketplace/actions/git-town-github-action)
to automatically add and update breadcrumbs on all pull requests. This offloads
the update workload to CI and ensures breadcrumbs get updated even when changes
are made outside the local Git Town CLI.

To enable this behavior, set up the GitHub Action.
