# Stacked changes

[Stacked changes](https://newsletter.pragmaticengineer.com/p/stacked-diffs) let
you implement and review complex work as a series of smaller, focused feature
branches that build on top of each other.

Key benefits of stacked changes:

- Developers and reviewers maintain momentum and block each other less
- Large, complex changes are broken into smaller, easier-to-review parts
- Merge conflicts are reduced by shipping already approved parts separately from
  work still under review

Implementing a complex change as a stack of branches requires running a lot more
Git commands. Git Town provides first-class support for stacked changes and
automates this extra work for you.

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
together while reviewing them independently. This is a perfect use case for
stacked branches.

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
renames, not the refactor diff.

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
git commit --down=2
```

This command does the following things:

1. Commit the currently staged changes into `1-refactor`
2. Pulls updates from `1-refactor` into `2-rename-foo`
3. Pulls updates from `2-rename-foo` into `3-rename-bar`

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

[git town sync](commands/sync.md) updates all branches:

1. Pulls updates made by other people into our local `main` branch
2. Deletes branch `1-refactor` from our local Git workspace because it was
   shipped at the remote
3. Pulls updates from `main` into `2-rename-foo`
4. Pulls updates from `2-rename-foo` into `3-rename-bar`

## Build the new feature

We can now add the new feature on top of the cleaned-up code base:

```
git town append 4-add-feature
```

Now you have a clean, reviewable stack:

- Each change lives in its own branch
- Branches build on top of each other
- All branches get reviewed independently
- [git town hack](https://www.git-town.com/commands/hack.html) starts a stack by
  creating its first branch
- [git town append](https://www.git-town.com/commands/append.html) extends a
  stack by adding a branch to its end
- You always ship the oldest branch in the stack
- `git town sync` keeps the stack up to date with other changes made to the
  codebase

## Best practices

### One change per branch

The
[single responsibility principle](https://en.wikipedia.org/wiki/Single-responsibility_principle)
applies to feature branches just as it does to functions, classes, and modules.
Single-responsibility branches are easier to reason about, less likely to
conflict, and allow shipping work faster.

When you have an idea that is different from what you currently work on, resist
the urge to code it in the current feature branch. Implement it in its own
feature, parent, or child branch.

If you can't create a new branch right now, write your idea down and implement
it later.

### Avoid unnecessary stacking

Only stack changes that depend on each other. If they don't, create them as
independent top-level feature branches that have `main` as their parent. This
setup has the advantage that you can ship any branch in any order.

It's okay to have multiple stacks.

### Keep your stack organized

Branches must be shipped oldest-first. Git Town provides powerful commands to
organize the branches in your stack:

- [git town hack](commands/hack.md) starts a new stack
- [git town append](commands/append.md) appends a new branch to the end of a
  stack
- [git town prepend](commands/prepend.md) inserts a new branch between the
  current branch and its parent
- [git town detach](commands/detach.md) extracts a branch from a stack and makes
  its own independent stack
- [git town swap](commands/swap.md) switches the position of the current branch
  and its parent in the stack
- [git town set-parent](commands/set-parent.md) changes the parent for the
  current branch and all its descendents

### Navigate your stack efficiently

To help commit the right changes to the right branch, Git Town provides powerful
commands to navigate stacks:

- [git town branch](commands/branch.md) shows you where you are in the stack
  hierarchy
- [git town switch](commands/switch.md) allows you to jump to any branch using a
  visual dialog with VIM motions
- [git town down](commands/up.md) switches to the parent branch
- [git town up](commands/up.md) switches to the child branch
- [git town walk](commands/walk.md) executes a CLI command or opens an
  interactive shell on each branch of the stack

### Embed the stack lineage into pull requests

The
[Git Town GitHub Action](https://github.com/marketplace/actions/git-town-github-action)
adds a visual graph of which branch of the stack the pull request is for. This
provides context when reviewing changes.

### Keep the stack in sync

Stacks are more prone to phantom merge conflicts than stand-alone branches. Run
`git town sync --stack` or `git town sync --all` regularly to propagate changes
across your stacks.

### Avoid phantom conflicts

_Phantom conflicts_ occur when Git reports a merge or rebase conflict that isn't
a real conflict. They can occur when multiple branches in a stack modify the
same line in the same file, and you ship using squash-merges.

After shipping the oldest branch from a stack using a squash-merge, `main`
contains a new commit with the same changes as the shipped branch but a
different commit hash. When syncing, Git sees the new commit on main and the
commit on the shipped branch as conflicting edits to the same line.

Git Town can detect and automatically resolve many of these phantom conflicts
because it tracks the branch hierarchy and understands the relationships between
commits.

To minimize phantom conflicts:

1. **Sync frequently.** In a synced stack, each branch builds directly on top of
   its parent, so changes are linear and easy for Git to reconcile. Branches in
   an unsynced stack drift apart, making conflicts more likely.

   If syncing takes too long, use
   [--detached](commands/sync.md#-d--detached--no-detached) and
   [--no-push](commands/sync.md#--push--no-push) to speed it up.

2. **Enable [rerere](https://git-scm.com/book/en/v2/Git-Tools-Rerere).** Git
   remembers how you resolved past conflicts and reuses those resolutions
   automatically.

3. **Ship using
   [fast-forward merges](https://git-scm.com/docs/git-merge#_fast_forward_merge).**
   Fast-forwarding keeps commit history between your stack and `main` identical,
   avoiding synthetic differences that cause phantom conflicts.

   - [GitLab supports this natively](https://docs.gitlab.com/ee/user/project/merge_requests/methods/#fast-forward-merge).
   - On GitHub, use [git town ship](commands/ship.md) with the
     [fast-forward strategy](preferences/ship-strategy.md#fast-forward) to
     achieve the same effect. See GitHub’s
     [docs](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/about-pull-request-merges#squashing-and-merging-a-long-running-branch)
     for details.

4. **Compress noisy branches.** If a branch has too many commits and keeps
   hitting the same conflicts, [compress](commands/compress.md) it to a single
   commit.

5. **Keep branches focused.** Small, single-purpose branches make it easier to
   understand and resolve conflicts, and to see what changed, why, and where.
