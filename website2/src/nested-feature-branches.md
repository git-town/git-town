# Feature branch chains

The
[single responsibility principle](https://en.wikipedia.org/wiki/Single-responsibility_principle)
applies to feature branches the same way it applies to code architecture.
Implementing, refactoring, reviewing, and resolving merge conflicts are much
more straightforward on a branch that implements a single change than on a
branch that contains several unrelated changes.

Git Town's _branch chain_ feature supports working with single-responsibility
feature branches. You implement each logical code change in a separate feature
branch organize all branches in a chain that makes changes made in parent
branches visible to their child branches.

As an example, let's say we want to add a feature to an existing codebase.
Before can do that cleanly, we need to get the code base ready. In particular,
we need to:

1. Clean up some technical drift: improve the names of variables, functions,
   classes, and files whose meaning has changed over time as the application has
   evolved. Cleaning this up now allows us to use more accurate names for our
   new feature.
2. Make the architecture more flexible so that we can add the new feature in a
   clean way.
3. Build the feature on top of the new architecture made in (1) and (2).
4. While looking through the code we also found some typos that we want to fix.

Implementing all these changes in a single feature branch is a big and risky
change that will take a long time. That's problematic in several ways. Touching
so many files in the codebase and on a branch that exists until we have finished
all these features will cause substantial merge conflicts with changes made by
other team members in that time. If we mess up one of the changes and want to
start over, we would have to throw away all the other changes as well. If the
code review for any of the four features takes longer, it will hold back
shipping the other features. Hence, let's implement these changes in separate
branches. But since feature (3) depends on the changes in (1) and (2), and
drives the changes in (2), we want to develop them together. The solution is a
chain of feature branches.

## Branch 1: fix typos

First, let’s fix the typos because there is no reason to keep looking at them.
We create a feature branch named `1-fix-typos` to contain the typo fixes using
the [git hack](commands/hack.md) command we already know:

```
git hack 1-fix-typos
```

We fix the typos and submit a pull request via
[git new-pull-request](commands/new-pull-request.md).

This took under a minute. While these changes get reviewed, we move on to fixing
the technical drift.

## Branch 2: rename foo

We don’t want to look at the typos that we just fixed again, so let’s perform
any further changes on top of branch `1-fix-typos`:

```
git append 2-rename-foo
```

[git append](commands/append.md) creates a new feature branch on top of the
current branch (which is `1-fix-typos`). We now have this branch hierarchy:

```
main
  \
   1-fix-typos
     \
      2-rename-foo
```

Now we commit the changes that rename the foo variable and start the next pull
request.

Because we used `git append` to create the new branch, Git Town knows about the
branch hierarchy and creates a pull request from branch `2-rename-foo` against
branch `1-fix-typos`. This guarantees that the pull request for branch 2 shows
only the changes made in that branch (renaming the variable) and not the syntax
fixes made in branch 1.

## Branch 3: rename bar

This is a different change from renaming `foo` and have different reviewers.
Let's do it in a different branch. Some of these changes might happen on the
same lines where we also renamed `foo` earlier. We don't want to have to deal
with merge conflicts coming from that. So let's make this change on top of the
change we made in step 2:

```
git append 3-rename-bar
```

We end up with this branch hierarchy:

```
main
  \
   1-fix-typos
     \
      2-rename-foo
        \
         3-rename-bar
```

## Fixing more typos

While renaming `bar`, we discovered more typos. Let's add them to the first
branch.

```
git checkout 1-fix-typos
# make the changes and commit them
git checkout 3-rename-bar
```

Back on branch `3-rename-bar`, the freshly fixed typos are visible again because
the commit to fix them exists only in branch `1-fix-typos` right now. Let's
propagate these changes through the entire branch chain so that they become
visible in branches 2 and 3 as well:

```
git sync
```

## Branch 4: generalize the infrastructure

With everything appropriately named we can make larger changes. We cut branch
`4-generalize-infrastructure` and perform the refactor in it. It has to be a
child branch of `3-rename-bar`, since the improved variable names done in the
latter branch will help make the larger changes we are about to do now.

```
git append 4-generalize-infrastructure
```

This refactoring touches a lot of files so we want to get this done and shipped
as fast as possible. Since that’s all we do in this branch, it’s pretty
straightforward to do and review. Off goes the code review for those changes.

## Shipping the typo fixes

In the meantime, we got the approval for the typo fixes in step 1. Let’s ship
them!

```
git ship 1-fix-typos
```

You have to use the [git ship](commands/ship.md) command here because it updates
the branch hierarchy that Git Town keeps track of.

With branch `1-fix-typos` shipped, our branch hierarchy now looks like this:

```
main
  \
   2-rename-foo
    \
     3-rename-bar
       \
        4-generalize-infrastructure
```

## Synchronizing our work with the rest of the world

We have been at it for a while. Other developers on the team have shipped
branches too. We don't want our branches to deviate too much from what’s
happening on the `main` branch since that can lead to more severe merge
conflicts later. Let's get everything in sync!

```
git sync
```

This merges `main` into `2-rename-foo`, then `2-rename-foo` into `3-rename-bar`,
then `3-rename-bar` into `4-generalize-infrastructure`.

## Branch 5: building the new feature

With the new generalized code architecture in place, we can now add the new
feature.

```
git append 5-add-feature
```

Let’s stop here and review what we have done.

- Each change happens in its own feature branch.
- All feature branches form a chain that makes changes made in earlier branches
  visible in ones later in the chain
- Each feature branch in the chain gets reviewed and shipped in isolation.
- `git hack` creates a feature branch as a child of the main branch.
- `git append` creates a new feature branch as a child of the current feature
  branch. This starts a feature branch chain.
- `git sync` keeps all feature branches in sync with the rest of the world.
- `git ship` ships the oldest feature branch in a branch chain.

## Advantages

There are many advantages of implementing a large code change as a chain of
feature branches:

- Changes are more focused: each branch makes a single change which is easier to
  reason about and quicker to implement, debug, review, and ship than a more
  complex change.
- Branches containing focused changes cause less and smaller merge conflicts
  that are easier to resolve than branches that contain many different changes.
- You can start the review/ship reviewing parts of your changes sooner because
  you start submitting pull requests earlier.

Ultimately, using this technique you will get more work done faster. You have
more fun because there is a lot less getting stuck, spinning wheels, and
starting over. Working this way requires running a lot more Git commands, but
with Git Town this is a complete non-issue since it automates this repetition
for you. Best Practices

To fully leverage this technique, all you have to do is follow a few simple
rules:

postpone ideas: when you work on something and run across an idea for another
change, resist the urge to do it right away. Instead, write down the change you
want to do (on a sheet of paper or a simple text file), finish the change you
are working on right now, and then perform the new change in a new branch a few
minutes later. If you can’t wait at all, commit your open changes into the
current branch, create the next branch, perform the new changes there, then
return to the previous branch and finish the work there.

go with one chain of branches: When in doubt whether changes depend on previous
changes and might or might not cause merge conflicts later, just work in child
branches. It has almost no side effects, except that you have to ship the
ancestor branches first. If your branches are focused, you will get very fast
reviews, be able to ship them quickly, and they won’t accumulate.

do large refactorings first: In our example, we did the refactor relatively late
in the chain because it wasn’t that substantial. Large refactorings that touch a
lot of files have the biggest potential for merge conflicts with changes from
other people, though. You don’t want them hanging around for too long, but get
them shipped as quickly as you can. You can use git prepend to insert feature
branches before the currently checked out feature branch. If you already have a
long chain of unreviewed feature branches, try to insert the large refactor into
the beginning of your chain, so that it can be reviewed and shipped as quickly
as possible:

git checkout 2-rename-foo git prepend 1-other-large-refactor

This leads to the following branch hierarchy:

master\
1-other-large-refactor\
2-rename-foo\
3-rename-bar\
4-generalize-infrastructure

The new large refactor is at the front of the line, can be shipped right when it
is reviewed, and our other changes are now built on top of it.

Happy hacking!
