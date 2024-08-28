# ship-strategy

The `ship-strategy` setting defines how [git town ship](../commands/ship.md)
merges finished feature branches into the main branch.

## options

### api

When using the "api" ship-strategy, [git ship](../commands/ship.md) presses the
"merge" button for the proposal in the web UI of your code hosting platform via
an API call.

You need to configure an API token in the
[setup assistant](../commands/config-setup.md) for this to work.

`api` is the default value because it does exactly what you normally do
manually.

### squash-merge

When set to `squash-merge`, [git ship](../commands/ship.md) merges the feature
branch to ship on your local Git repository.

feature branches against their parent branches and then does a safe force-push
of your rebased local commits to the tracking branch. This safe force-push uses
Git's
[--force-with-lease](https://git-scm.com/docs/git-push#Documentation/git-push.txt---no-force-with-lease)
and
[--force-if-includes](https://git-scm.com/docs/git-push#Documentation/git-push.txt---no-force-if-includes)
switches to guarantee that the force-push will never overwrite commits on the
tracking branch that haven't been integrated into the local Git history.

If the safe force-push fails, Git Town rebases your local branch against its
tracking branch to pull in new commits from the tracking branch. If that leads
to conflicts, you have a chance to resolve them and continue syncing by running
[git town continue](../commands/continue.md).

When continuing the sync this way, Git Town tries again to safe-force-push and
rebase until the safe-force-push succeeds without removing commits from the
tracking branch that aren't part of the local Git history.

This can lead to an infinite loop if you do an interactive rebase that removes
commits from the tracking branch while syncing it. You can break out of this
infinite loop by doing a less aggressive rebase that doesn't remove the remote
commits. Finish the `git sync` command and then clean up your commits via a
separate interactive rebase after the sync. At this point another sync will
succeed because the commits you have just cleaned up are now a part of your
local Git history.

The rule of thumb is that pulling in new commits via `git sync` and cleaning up
old commits must happen separately from each other. Only then can Git guarantee
that the necessary force-push happens without losing commits.

### compress

When using the `compress` sync strategy, [git sync](../commands/sync.md) first
merges the tracking and parent branches and then
[compresses](../commands/compress.md) the synced branch.

This sync strategy is useful when you want all your pull requests to always
consists of only one commit.

Please be aware that this sync strategy leads to more merge conflicts than the
"merge" sync strategy when more than one Git user makes commits to the same
branch.

## change this setting

The best way to change this setting is via the
[setup assistant](../configuration.md).

### config file

In the [config file](../configuration-file.md) the sync-feature-strategy is part
of the `[sync-strategy]` section:

```toml
[sync-strategy]
feature-branches = "merge"
```

### Git metadata

To manually configure the sync-feature-strategy in Git, run this command:

```
git config [--global] git-town.sync-feature-strategy <merge|rebase>
```

The optional `--global` flag applies this setting to all Git repositories on
your local machine. When not present, the setting applies to the current repo.
