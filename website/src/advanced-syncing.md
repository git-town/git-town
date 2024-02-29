# Advanced Branch Syncing (beta)

Git branches can be used in many different ways. If you tell Git Town how you
use the branches in your workspace, you can always run `git sync` or
`git sync --all` and all your local branches will get synced the way you want
them to be synced.

## Branch Ownership

Branches that you create using [git hack](commands/hack.md),
[git append](commands/append.md), and [git prepend](commands/prepend.md) are
owned by you. Owning a branch means that you are the person that manages this
branch. You [propose](commands/propose.md), [ship](commands/ship.md), pull in
updates from its parent branch, and [delete](commands/kill.md) it from the
hosting server.

You might also have branches on your machine that you don't own. An example is
when you help somebody with a problem on a branch that this person owns. You
might want to do some experiments on copy of the other person's branch on your
machine. But you want to commit these experiments on your machine so that you
can switch branches. You don't want to pull in parent updates or accidentally
ship or delete the other person's branch from the hosting service because those
are responsibilities of the person owning that branch or feature.

## Observed branches

If you want that `git sync` doesn't push your local commits to the tracking
branch,

## Contribution branches

If you want that `git sync` pushes your local commits to the tracking branch,
you make this branch a `contribution branch`.

## Parked Branches

Parked branches don't get synced unless you run `git sync` directly on a parked
branch. You might want to park a branch if you

- want to intentionally keep the branch at an older state
- don't want to deal with merge conflicts on this branch right now
- reduce load on your CI server
