# Ship strategy

This setting defines how [git town ship](../commands/ship.md) merges finished
feature branches into the main branch.

## options

### api

When using the "api" ship strategy, [git town ship](../commands/ship.md) presses
the "merge" button for the proposal in the web UI of your forge via an API call.

You need to configure an API token in the [setup assistant](../commands/init.md)
for this to work.

`api` is the default value because it does exactly what you normally do
manually.

### always-merge

The `always-merge` ship strategy creates a merge commit via `git merge --no-ff`.

This strategy allows visually grouping related feature commits together which
may aid in understanding project history in certain situations.

It is not generally recommended to revert merge commits, so `git town undo` will
not create a merge reversal commit if the merge commit has been pushed already.
See
[howto/revert-a-faulty-merge.adoc](https://github.com/git/git/blob/master/Documentation/howto/revert-a-faulty-merge.adoc)
in the official Git documentation for more information.

### fast-forward

The `fast-forward` ship strategy prevents false merge conflicts when using
[stacked changes](../stacked-changes.md) and allows to
[Ship several branches in a stack](../how-to/ship-stack.md) without unnecessary
CI runs. It runs
[git merge --ff-only](https://git-scm.com/docs/git-merge#Documentation/git-merge.txt---ff-only)
which fast-forwards the parent branch to contain the commits of the branch to
ship and then pushes the new commits on the parent branch to the
[development remote](dev-remote.md). This way the parent branch contains the
exact same commits as the branch that has just been shipped.

For details why this is needed check out this
[GitHub documentation](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges#squashing-and-merging-a-long-running-branch).

This works on GitHub even if your main branch is protected as long as the
associated proposal is green and has been approved! GitHub recognizes that the
commits you push have already been tested and approved and allows them to be
pushed. For more information, see
[this StackOverflow answer](https://stackoverflow.com/questions/60597400/how-to-do-a-fast-forward-merge-on-github/66906599#66906599).

A limitation of the `fast-forward` ship strategy is that your feature branch
must be up to date, i.e. the main branch must not have received additional
commits since you last synced your feature branch.

### squash-merge

When set to `squash-merge`, [git town ship](../commands/ship.md) merges the
feature branch to ship in your local Git repository. While doing so it squashes
all commits on the feature branch into a single commit and lets you edit the
commit message.

### config file

Set the ship strategy in the [config file](../configuration-file.md):

```toml
[ship]
strategy = "api"
```

### Git metadata

To manually configure the ship strategy in Git metadata, run:

```wrap
git config [--global] git-town.ship-strategy <always-merge|api|fast-forward|squash-merge>
```

The optional `--global` flag applies this setting to all Git repositories on
your machine. Without it, the setting applies only to the current repository.

## environment variable

You can configure the ship strategy by setting the `GIT_TOWN_SHIP_STRATEGY`
environment variable.
