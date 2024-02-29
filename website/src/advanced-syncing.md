# Advanced Branch Syncing

Git branches can be used in many different ways. When Git Town is configured
correctly, you can run `git sync` or `git sync --all` any time and each of your
local branches will get synced in the specific ways it's supposed to get synced
or not synced.

## Branch Ownership

The "owner" of a branch is responsible for the larger lifetime events of this
branch. The branch owner pulls in updates from the parent branch, creates the
pull request, and delete the branch from the hosting server. Typically you own
branches that you create with [git hack](commands/hack.md),
[git append](commands/append.md), or [git prepend](commands/prepend.md).

You might also have branches on your machine that you don't own. An example is
when you review somebody else's branch in your local editor. Or when you help
somebody solve a problem that happens on a branch this person owns.

In both cases, you don't want to pull in parent updates or accidentally ship or
delete the other person's branch from the hosting service because those are
responsibilities of the person owning that branch. But you want your local
branch to receive additional commits made to this branch.

## Contribution branches

Contribution branches are for people who contribute commits to somebody else's
branch. `git sync` always [rebases](preferences/sync-feature-strategy#rebase)
your local commits on a contribution branch. It does not pull updates from the
parent branch. You cannot [propose](commands/propose.md) or
[ship](commands/ship.md) contribution branches. `git sync` removes contribution
branches from your machine as soon as their tracking branch is gone, even if you
have unpushed local commits. When you [kill](commands/kill.md) a contribution
branch, it only deletes your local copy and not the tracking branch.

Run [git contribute](commands/contribute.md) on a branch to make it a
contribution branch.

## Observed branches

Observed branches are for people who want to observe the work of somebody else
without contributing to it. `git sync` only pulls updates from the tracking
branch, always using the [rebase](preferences/sync-feature-strategy#rebase)
sync-feature-strategy. It doesn't push your local commits. You cannot
[propose](commands/propose.md) or [ship](commands/ship.md) observed branches.
`git sync` removes observed branches from your machine as soon as their tracking
branch is gone, even if you have unpushed local commits. When you
[kill](commands/kill.md) an observed branch, it only deletes your local copy and
not the tracking branch.

Run [git observe](commands/observe.md) on a branch to make it an observed
branch.

## Parked Branches

Parked branches don't get synced at all unless you run `git sync` directly on a
parked branch. You might want to park a branch if you

- want to intentionally keep the branch at an older state
- don't want to deal with merge conflicts on this branch right now
- reduce load on your CI server by syncing only your actively developed local
  branches

Run [git park](commands/park.md) on a branch to park it. You cannot park the
perennial branches including the main branch.
