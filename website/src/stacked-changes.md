# Stacked changes

[Stacked changes](https://newsletter.pragmaticengineer.com/p/stacked-diffs) let
you implement and review complex work as a series of smaller, focused feature
branches that build on top of each other.

Key benefits of stacked changes:

- Developers and reviewers maintain momentum and block each other less
- Large, complex changes are broken into smaller, easier-to-review parts
- Merge conflicts are reduced by shipping already approved parts separately from
  work still under review

Git Town provides first-class support for stacked changes.

## Example

Suppose you are adding a new feature to an existing codebase. Before we can do
that cleanly, you need to prepare the code base:

1. Refactor the architecture to make it easier to add the new feature cleanly
2. Clean up technical drift: rename variables, functions, etc
3. Build the feature on top of the improved codebase

Putting all these changes into one feature branch is risky. The refactor in (1)
touches many files that other people may also be changing. We want to review and
merge this quickly to minimize conflicts. The feature in (3) might take longer
to build. Both changes should therefore live in separate branches.

However, the feature in (3) depends on (1) and (2). We need to develop them
together while reviewing them independently. The perfect use case for stacked
branches.

## Branch 1: refactor

Start by creating a branch for the refactor:

```
git town hack 1-refactor
```

[git town hack](commands/hack.md) creates a new feature branch off the `main`
branch. Implement the refactor and commit your changes.

## Branch 2: rename foo

Next, perform some renames that depend on the refactor. Create a new branch on
top of `1-refactor`:

```
git town append 2-rename-foo
```

[git town append](commands/append.md) creates a new feature branch on top of the
current branch. The resulting lineage looks like this:

```
main
 \
  1-refactor
   \
*   2-rename-foo
```

Branch `2-rename-foo` now includes the refactor from branch 1. When you open a
PR, Git Town will target `1-refactor` automatically, so reviewers see only the
renames — not the refactor diff.

## Branch 3: rename bar

This change is independent of renaming `foo` and may have a different reviewer.
Create another branch on top of `2-rename-foo`:

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

While working on `3-rename-bar`, you discover another improvement for the
architecture. Add it to `1-refactor`:

```sh
git checkout 1-refactor
# make the changes and commit them
git checkout 3-rename-bar
```

Your new refactor changes exist only in `1-refactor`. To propagate them through
the other branches in the stack, run:

```
git town sync --stack --detached
```

This:

1. Pulls updates from `1-refactor` into `2-rename-foo`
2. Pulls updates from `2-rename-foo` into `3-rename-bar`

## Shipping the refactor

Once the refactor is approved, you or somebody else merges this pull request.
The stack now looks like this:

```
main
 \
  1-refactor (the remote branch is gone, the local branch still exists)
   \
    2-rename-foo
     \
*     3-rename-bar
```

## Keeping the stack up to date

We have been at it for a while. Other team members made changes to the codebase
as well. We don't want our local branches to deviate too much from the rest of
the codebase, since that leads to merge conflicts later. Let's get our local Git
workspace in sync with the rest of the universe!

```
git town sync --all
```

[git town sync](commands/sync.md) updates all branches in order:

1. Pulls updates made by other people into our local `main` branch
2. Deletes branch `1-refactor` from our local Git workspace because it was
   shipped at the remote
3. Pulls updates from `main` into `2-rename-foo`
4. Pulls updates from `2-rename-foo` into `3-rename-bar`

## Build the new feature

We can now add the new feature on top of the cleanedl-up code base:

```
git town append 4-add-feature
```

Now you have a clean, reviewable stack:

- Each change lives in its own branch
- Branches build on top of each other
- All branches get be reviewed independently
- [git town hack](https://www.git-town.com/commands/hack.html) starts a stack by
  creating its first branch
- [git town append](https://www.git-town.com/commands/append.html) extends a
  stack by adding a branches to its end
- You always ship the oldest branch in the stack
- `git town sync` keeps the stacks current

Single-responsibility branches are easier to reason about, less likely to
conflict, and allow shipping work faster. Implementing a complex change as a
stack of branches requires running more Git commands. Git Town automates this
extra work for you.

## Best practices

#### One change per branch

The
[single responsibility principle](https://en.wikipedia.org/wiki/Single-responsibility_principle)
applies to feature branches just as it does to functions, classes, and modules.
Each branch should only make a single, consistent change. Such single-purpose
branches are easier to implement, refactor, review, test and merge than branches
that mix unrelated changes.

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

#### Minimize commit changes when shipping

When using stacks, try to fast-forward your feature branches into the main
branch to avoid empty merge conflicts when syncing the stack later. On GitLab
that's
[straightforward](https://docs.gitlab.com/ee/user/project/merge_requests/methods/#fast-forward-merge).
GitHub does not provide a fast-forward merge option out of the box but you can
achieve it with the
[fast-forward ship strategy](preferences/ship-strategy.md#fast-forward) together
with the [compress](preferences/sync-feature-strategy.md#compress) or
[rebase](preferences/sync-feature-strategy.md#rebase) sync strategy. The
[Git Town GitHub Action](https://github.com/marketplace/actions/git-town-github-action)
adds a visual description of which branch of the stack the pull request is for.

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
   [--no-push](commands/sync.md#--push--no-push) flags to speed it up.

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
