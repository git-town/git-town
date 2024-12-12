# Branch Types

Git Town supports many different types of Git branches. When properly
configured, you can run `git town sync` or `git town sync --all` at any time and
each of your local branches will get synced in the specific ways it's supposed
to get synced or not synced.

## Feature branches

Feature branches are the branches on which you typically make changes. They are
typically cut from the _main branch_ and get merged back into it. You can also
cut feature branches from any other branch type if needed. Feature branches sync
with their parent and tracking branch.

## Main branch

The main branch is a _perennial branch_ from which feature branches get cut by
default. The main branch contains the latest development version of your
codebase.

## Perennial branches

Perennial branches are long-lived branches. They have no parent and are never
shipped. Typical perennial branches are `main`, `master`, `development`,
`production`, `staging`, etc. Perennial branches often correspond with a cloud
environment of the same name.

## Contribution branches

Contribution branches are for people who contribute commits to somebody else's
feature branch. You cannot [propose](commands/propose.md) or
[ship](commands/ship.md) contribution branches because those are
responsibilities of the person owning the branch you contribute to. For the same
reason `git town sync` does not pull updates from the parent branch of a
contribution branch and always
[rebases](preferences/sync-feature-strategy.md#rebase) your local commits.
Syncing removes contribution branches from your machine as soon as their
tracking branch is gone, even if you have unpushed local commits.
[Deleting](commands/delete.md) a contribution branch only deletes your local
copy and not the tracking branch.

You can make any feature branch a contribution branch by running
[git town contribute](commands/contribute.md) on it. Convert a contribution
branch back to a feature branch by running [git town hack](commands/hack.md) on
it. You can also define a
[contribution-regex](preferences/contribution-regex.md) in your Git
configuration or the config file.

## Observed branches

Observed branches are for people who want to observe the work of somebody else
without contributing commits to it. Similar to contribution branches, you cannot
[propose](commands/propose.md) or [ship](commands/ship.md) observed branches,
[delete](commands/delete.md) only deletes your local copy and not the tracking
branch, `git town sync` always uses the
[rebase](preferences/sync-feature-strategy.md#rebase) sync-feature-strategy and
will remove a local observed branch as soon as its tracking branch is gone, even
if there are unmerged local commits.

Unlike with contributing branches, `git town sync` does not push your local
commits made to an observed branch to its tracking branch.

You can make any feature branch an observed branch by running
[git town observe](commands/observe.md) on it. Convert an observed branch back
to a feature branch by running [git town hack](commands/hack.md) on it. You can
also define an [observed-regex](preferences/observed-regex.md) in your Git
configuration or the config file.

## Parked Branches

Parked branches don't get synced at all unless you run `git town sync` directly
on a parked branch. You might want to park a branch if you

- want to intentionally keep the branch at an older state
- don't want to deal with merge conflicts on this branch right now
- reduce load on your CI server by syncing only your actively developed local
  branches

You can park any feature branch by running [git town park](commands/park.md) on
it. Unpark a parked branch by running `git town hack` on it.

## Prototype Branches

A prototype branch is a local-only feature branch that incorporates updates from
its parent branch but is not pushed to the remote repository. Prototype branches
are useful when:

- the branch contains sensitive information, such as secrets, or potentially
  problematic code or data that could trigger alerts
- the developer prefers to keep their work private from the rest of the team
  during the initial stages of development
- you want to reduce CI pressure in the early phases of feature development when
  there isn't anything to test

Syncing prototype branches follows the
[sync-prototype-strategy](preferences/sync-prototype-strategy.md) or - if this
setting isn't present - the
[sync-feature-strategy](preferences/sync-feature-strategy.md). This allows you
to rebase your commits while working locally, and avoid rebasing when your
commits become visible to others.

When you [propose](commands/propose.md) a prototype branch, it loses its
prototype status since it now has an official tracking branch that other people
look at. In this situation you can keep syncing without pushes by using the
`--no-push` sync option.

You can compress and ship prototype branches as usual. Parking and unparking a
prototype branch maintains its prototype status. When you change a prototype
branch to an observed or contribution branch it loses its prototype status.

To designate any feature branch as a prototype branch, execute
[git town prototype](commands/prototype.md) on it. To convert a prototype branch
to a feature branch, use [git town hack](commands/hack.md).

## Configuring branch types

You can set the types of indivdiual branches with these commands:

- [contribute](commands/contribute.md)
- [hack](commands/hack.md)
- [observe](commands/observe.md)
- [park](commands/park.md)
- [prototype](commands/prototype.md)

These preferences allow you to configure the types of larger groups of branches:

- [default-branch-type](preferences/default-branch-type.md),
- [feature-regex](preferences/feature-regex.md), and
- [new-branch-type](preferences/new-branch-type.md),
- [perennial-regex](preferences/perennial-regex.md) preferences.
