# Advanced Branch Syncing (beta)

Git branches can be used in many different ways. If you tell Git Town how you
use the branches in your workspace, you can always run `git sync` or
`git sync --all` and all your local branches will get synced the way you want
to.

## Parked Branches

## Branch Ownership

Branches that you create using [git hack](commands/hack.md),
[git append](commands/append.md), and [git prepend](commands/prepend.md) are
owned by you. Owning a branch means that you are the person that manages this
branch. You [propose](commands/propose.md), [ship](commands/ship.md), pull in
updates from its parent branch, or [delete](commands/kill.md) it.

You might also have branches on your machine that you don't own. An example is
when you help another person with a problem on a branch that this person owns.
You might want to do some experiments on copy of the other person's branch on
your machine. You want to commit these experiments on your machine so that you
can switch branches.

## Contribution branches

If you
