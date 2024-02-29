# Advanced Branch Syncing

Git branches can be used in many different ways. When properly configured, you
can run `git sync` or `git sync --all` at any time and each of your local
branches will get synced in the specific ways it's supposed to get synced or not
synced.

## Contribution branches

Contribution branches are for people who contribute commits to somebody else's
branch. You cannot [propose](commands/propose.md) or [ship](commands/ship.md)
contribution branches because those are responsibilities of the person owning
the branch you contribute to. For the same reason `git sync` does not pull
updates from the parent branch of a contribution branch and always
[rebases](preferences/sync-feature-strategy#rebase) your local commits. Syncing
removes contribution branches from your machine as soon as their tracking branch
is gone, even if you have unpushed local commits. [Killing](commands/kill.md) a
contribution branch only deletes your local copy and not the tracking branch.

You can make any feature branch a contribution branch by running
[git contribute](commands/contribute.md) on it. Convert a contribution branch
back to a feature branch by running [git hack](commands/hack.md) on it.

## Observed branches

Observed branches are for people who want to observe the work of somebody else
without contributing commits to it. Like contributing branches, you cannot
[propose](commands/propose.md) or [ship](commands/ship.md) observed branches,
[kill](commands/kill.md) only deletes your local copy and not the tracking
branch, `git sync` always uses the
[rebase](preferences/sync-feature-strategy#rebase) sync-feature-strategy and
will remove a local observed branch as soon as its tracking branch is gone, even
if there are unmerged local commits.

Unlike with contributing branches, `git sync` does not push your local commits
made to an observed branch to its tracking branch.

You can make any feature branch an observed branch by running
[git observe](commands/observe.md) on it. Convert an observed branch back to a
feature branch by running [git hack](commands/hack.md) on it.

## Parked Branches

Parked branches don't get synced at all unless you run `git sync` directly on a
parked branch. You might want to park a branch if you

- want to intentionally keep the branch at an older state
- don't want to deal with merge conflicts on this branch right now
- reduce load on your CI server by syncing only your actively developed local
  branches

You can park any feature branch by running [git park](commands/park.md) on it.
You cannot park the perennial branches as well as the main branch. Unpark a
parked branch by running `git hack` on it.
