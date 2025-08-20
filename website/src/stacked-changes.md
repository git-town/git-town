# Stacked changes

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

We got the approval for the refactor from step 1. Let's ship it!

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

If you ship feature branches through the API or UI of your forge, run
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

Let's stop here and review what we have done.

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

## Best practices

#### One change per branch

When you have an idea that is different from what you currently work on, resist
the urge to code it in the current feature branch. Implement it in its own
feature, parent, or child branch.

If you can't create a new branch right now, write the idea down and implement it
later.

#### Keep the stack in sync

Stacks are more susceptible to phantom merge conflicts than stand-alone
branches. Don't forget to populate changes across all branches in your stack by
running `git town sync --stack` or `git town sync --all`.

#### Avoid unnecessary stacking

To reduce merge conflicts, feature branches should not diverge too much from the
main development branch. Stacking multiple changes on top of each other
amplifies this divergence. Overly "tall" stacks are therefore an anti-pattern to
avoid. It's often better to work in independent top-level feature branches by
default, and only stack branches if the changes they contain really depend on
each other. This way you can get your changes reviewed and shipped concurrently
and in any order, i.e. faster and with fewer merge conflicts.

#### Organize branch chains in the order you want to ship

You always have to ship the oldest branch first. You can use
[git town prepend](commands/prepend.md) to insert a feature branch as a parent
of the current feature branch or [set parent](commands/set-parent.md) to change
the order of branches.

#### Avoid phantom conflicts

_Phantom conflicts_ occur when Git reports a merge or rebase conflict that -
when looked at with more context - isn't a real conflict. Phantom conflicts can
occur when multiple branches in a stack modify the same line in the same file,
and you ship using squash-merges.

After you ship the oldest branch of such a stack, the main branch contains a new
commit that makes the same changes as the shipped branch, but as a different
commit than the one(s) on the shipped branch. As this new commit populates
through the stack in the next sync, Git sees sees two changes to the same lines
and assumes a conflict.

Git Town can resolve these phantom conflicts because it tracks the branch
hierarchy, can investigate such conflicts, and execute multiple Git commands to
resolve them.

Here are some best practices to minimize phantom merge conflicts:

1. Sync frequently. In a synced stack, each branch builds directly on top of its
   parent, so changes are linear and easy for Git to reconcile. In an unsynced
   stack, sibling branches evolve concurrently, making conflicts more likely,
   especially when they touch the same files.

   If you are hesitant to sync because it takes too long, use the
   [--detached](commands/sync.md#-d--detached--no-detached) and
   [--no-push](commands/sync.md#--no-push) flags to speed it up.

2. Enable Git's [rerere](https://git-scm.com/book/en/v2/Git-Tools-Rerere)
   feature. This lets Git remember how you resolved past conflicts and applies
   those resolutions in the future.

3. Ship using a
   [fast-forward merge](https://git-scm.com/docs/git-merge#_fast_forward_merge).
   This ensures the commits on main are byte-for-byte identical to those on the
   shipped branchs. This preserved shared history avoids unnecessary merges or
   rebases that likely produce phantom conflicts.

   - GitLab supports this
     [natively](https://docs.gitlab.com/ee/user/project/merge_requests/methods/#fast-forward-merge).
   - GitHub doesn’t support fast-forward merges via the UI, but you can achieve
     the same effect by [shipping locally](commands/ship.md) with Git Town's
     [fast-forward strategy](preferences/ship-strategy.md#fast-forward) and then
     pushing the result. See GitHub’s
     [docs](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges#squashing-and-merging-a-long-running-branch)
     for details.

4. If a feature branch has too many commits and you're resolving the same
   conflicts repeatedly, [compress](commands/compress.md) it down to a single
   commit.

5. Focus your feature branches to implement only a single change. This reduces
   the amount of context you need to process when resolving merge conflicts, and
   makes it easier to see which branch makes which change and why, and what the
   correct resolution is.
