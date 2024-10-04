# Stacked Changes

[Stacked changes](https://newsletter.pragmaticengineer.com/p/stacked-diffs)
implement and review a complex change as a series of smaller feature branches
that build on top of each other. Benefits of stacked changes are:

- developer and reviewer maintain momentum and block less on each other
- breaking up the problem of developing/reviewing a complex problem into
  developing/reviewing many smaller problems
- minimize merge conflicts by shipping parts of a complex change that are
  already approved separately from parts still under review

The
[single responsibility principle](https://en.wikipedia.org/wiki/Single-responsibility_principle)
applies to feature branches the same way it applies to functions, classes, and
methods. Feature branches should also perform only one change. Implementing,
refactoring, reviewing, and resolving merge conflicts on such
single-responsibility branches is easier than with branches that combine
unrelated changes.

Git Town provides wide reaching support for stacked changes. When using stacked
changes, try to fast-forward your feature branches into the main branch to avoid
empty merge conflicts when syncing the stack later. On GitLab that's
[straightforward](https://docs.gitlab.com/ee/user/project/merge_requests/methods/#fast-forward-merge).
GitHub does not provide a fast-forward merge option out of the box but you can
achieve it with the
[fast-forward ship strategy](preferences/ship-strategy.md#fast-forward) together
with the [compress](preferences/sync-feature-strategy.md#compress) or
[rebase](preferences/sync-feature-strategy.md#rebase) sync strategy. The
[Git Town GitHub Action](https://github.com/marketplace/actions/git-town-github-action)
adds a visual description of which branch of the stack the pull request is for.

## Example

Let's say we want to add a new feature to an existing codebase. Before we can do
that cleanly, we need to get the code base ready:

1. Make the architecture more flexible so that we can add the new feature in a
   clean way.
2. Clean up some technical drift: improve the names of variables and functions.
3. Build the feature on top of this modernized codebase

Implementing all these changes in a single feature branch is risky. Some changes
like the refactor in (1) touch a lot of files that other people might change as
well. We want to review and merge them as fast as possible to minimize merge
conflicts. Other changes like building the actual feature in (3) will take a
while to build. We should therefore make both changes in separate branches. At
the same time, the feature (3) depends on the changes in (1) and (2) and drives
the changes in (2). We want to develop these changes together. The solution is a
stack of feature branches.

## Branch 1: refactor

The first feature branch contains the refactor. We create a feature branch named
`1-refactor` off the main branch to contain it.

```
git town hack 1-refactor
```

[git town hack](commands/hack.md) creates a new feature branch off the main
branch. We perform the refactor and commit it.

## Branch 2: rename foo

With the refactored architecture in place, we can update the names of some
variables, functions, and files whose role has changed as the code base has
evolved. Since these changes require the refactor, we perform them on top of
branch `1-refactor`:

```
git town append 2-rename-foo
```

[git town append](commands/append.md) creates a new feature branch on top of the
currently checked out branch (which is `1-refactor`). We now have this lineage:

```
main
 \
  1-refactor
   \
*   2-rename-foo
```

Branch `2-rename-foo` builds on top of `1-refactor` and thereby contains all the
changes made there. We commit the changes that rename the `foo` variable.
Because we used `git town append` to create the new branch, Git Town knows about
the lineage and creates the proposal (aka pull request) for branch
`2-rename-foo` against branch `1-refactor`. This way, the proposal for branch
`2-rename-foo` shows only the changes made in that branch (renaming the
variable) and not the refactor made in branch 1.

## Branch 3: rename bar

This is a different change from renaming `foo` and has different reviewers.
Let's perform it in a different branch. Some of these changes might happen on
the same lines on which we also renamed `foo` earlier. We don't want to deal
with merge conflicts coming from that. So let's make this change on top of the
change we made in step 2:

```
git town append 3-rename-bar
```

The lineage is now:

```
main
 \
  1-refactor
   \
    2-rename-foo
     \
*     3-rename-bar
```

## Extend the refactoring

While renaming `bar`, we discover another improvement for the architecture.
Let's add it to the refactoring branch.

```
git checkout 1-refactor
# make the changes and commit them
git checkout 3-rename-bar
```

Back on branch `3-rename-bar`, the additional refactor we just added isn't
visible because the commit for it exists only in branch `1-refactor` right now.
Let's propagate these changes through the entire branch chain so that they
become visible in branches 2 and 3 as well:

```
git town sync
```

Because we created the branches with `git town append`, Git Town knows about the
branch lineage and [git town sync](commands/sync.md) can update all branches in
the right order. It updates the `main` branch, merges `main` into branch 1. Then
it merges branch 1 into branch 2 and branch 2 into branch 3.

## Shipping the refactor

We got the approval for the refactor from step 1. Let’s ship it!

```
git town ship 1-refactor
```

You have to use the [git town ship](commands/ship.md) command here because it
updates the lineage that Git Town keeps track of. With branch `1-refactor`
shipped, our lineage now looks like this:

```
main
 \
  2-rename-foo
   \
*   3-rename-bar
```

If you ship feature branches via the code hosting API or web UI, run
`git town sync --all`, or `git town sync` on the youngest child branch, to
update the lineage.

## Synchronizing our work with the rest of the world

We have been at it for a while. Other developers on the team have made changes
to the codebase as well. We don't want our branches to deviate too much from the
`main` branch since that leads to more severe merge conflicts later. Let's get
all our branches in sync with the rest of the world!

```
git town sync --all
```

This pulls updates for the `main` branch, then merges it into `2-rename-foo`,
then `2-rename-foo` into `3-rename-bar`.

## Branch 4: building the new feature

We can now add the new feature on top of the code base we prepared:

```
git town append 4-add-feature
```

Let’s stop here and review what we have done.

- Each change happens in its own feature branch.
- Our feature branches build on top of each other and see changes in their
  parent branches.
- We review and ship each feature branch in the chain in isolation.
- `git town hack` creates a feature branch as a child of the main branch.
- `git town append` creates a feature branch as a child of the current feature
  branch.
- `git town sync` keeps a feature branch chain up to date with the rest of the
  world
- `git town ship` ships the oldest feature branch in a branch chain.

Single-responsibility branches are easier to reason about and faster to
implement, debug, review, and ship than branches performing multiple changes.
They encounter fewer and smaller merge conflicts which are easier to resolve
than merge conflicts on branches that implement many different changes. You can
review and ship parts of your complex change before the entire change is
finished. You can still make different changes in parallel, just commit them to
the correct branch.

## Best Practices

_Branch discipline:_ when you have an idea that is different from what you
currently work on, resist the urge to code it in the current feature branch.
Implement it in its own feature, parent, or child branch.

_Keep the entire branch chain in sync:_ Make sure you run `git town sync --all`
or `git town sync` on the youngest child branch to keep the entire chain of
feature branches synced.

_Avoid unnecessary chaining:_ If your feature branches don't depend on each
other, put them in (independent) top-level feature branches. This way you can
ship them in any order.

_Organize branch chains in the order you want to ship:_ You always have to ship
the oldest branch first. You can use [git town prepend](commands/prepend.md) to
insert a feature branch as a parent of the current feature branch or
[set parent](commands/set-parent.md) to change the order of branches.

_Ship using fast-forward:_ Merge conflicts happen when you merge two branches
that change the same location in the same file. Git cannot (and should not)
decided which version to use and it lets the user sort it out.

to resolve with two different commits (commits with different SHAs) making the
same change. This

To minimize the risk of phantom merge conflicts after shipping the oldest branch
in a stacked change, ship using a
[fast-forward merging](https://git-scm.com/docs/git-merge#_fast_forward_merge).
This approach guarantees that the new commits on the main branch are the exact
same commits that were on the shipped feature branch. This helps Git resolve

GitLab provides this
[out of the box](https://docs.gitlab.com/ee/user/project/merge_requests/methods/#fast-forward-merge).
GitHub doesn't provide this out-of-the-box, but allows a workaround that you can
utilize by using [git town ship](commands/ship.md) with the
[fast-forward shipping strategy](preferences/ship-strategy#fast-forward). This
problem is documented by
[GitHub](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges#squashing-and-merging-a-long-running-branch).
