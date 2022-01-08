# Nested feature branches

The
[single responsibility principle](https://en.wikipedia.org/wiki/Single-responsibility_principle)
applies to feature branches as well. Implementing several different unrelated
code changes at the same time in the same feature branch is like trying to have
a conversation about several different topics at the same time with the same
person. It is never productive. We end up mixing issues or forgetting to think
about important edge cases of one topic because we are distracted with the other
topics. We only have parts of our brain available to think about each topic.
Dealing with several issues at the same time might create the illusion of
productivity, but in the end, it is faster, easier, safer, cleaner, and less
error-prone to take on each topic separately.

This blog post describes a technique for highly focused development by
implementing code changes using a series of nested Git branches. We use
specialized tooling (Git Town, an open-source plugin for Git) to make this type
of working easy and efficient. example

As an example, let’s say we want to implement a new feature for an existing
product. But to make such a complex change, we have to get the code base ready
for it first:

    clean up some technical drift by improving the names of classes and functions that don’t make sense anymore
    add some flexibility to the architecture so that the new feature can be built with less hackery
    while looking through the code base, we also found a few typos we want to fix

Let’s implement these things as a chain of independent but connected feature
branches! I provide the Git Town commands as well as the individual Git commands
you would have to run without Git Town for those unfamiliar with that tool.
First, let’s fix those typos because that’s the easiest change and there is no
reason to keep looking at them. fix typos

We create a feature branch named 1-fix-typos to contain the typo fixes from the
master branch:

git hack 1-fix-typos

git hack is a Git command added by Git Town. The corresponding vanilla Git
commands are:

git checkout master git pull git checkout -b 1-fix-typos

We always want to build new changes on top of the latest version of the master
branch.

We do a few commits fixing typos and submit a pull request:

git new-pull-request

This Git Town command opens a browser window to create the pull request on your
code hosting service.

All of this only took us under a minute. While the code review for those change
happens, we move on to fix the technical drift. rename foo

We don’t want to look at the typos we just fixed again, so let’s perform any
further changes on top of branch 1-fix-typos:

git append 2-rename-foo

git append creates a new feature branch by cutting it from the current branch,
resulting in this branch hierarchy:

master\
1-fix-typos\
2-rename-foo

The corresponding vanilla Git commands are:

git checkout -b 2-rename-foo

Now we commit the changes that rename the foo variable and start the next pull
request:

git new-pull-request

Because we used git append to create the new branch, Git Town knows about the
branch hierarchy and creates a pull request from branch 2-rename-foo against
branch 1-fix-typos. This guarantees that the pull request for branch 2 only
shows changes made in that branch (renaming the variable), but not the syntax
fixes made in branch 1. rename bar

This is a different change than renaming foo, so let's do it in a different
branch. Some of these changes might happen in the same places where we just
renamed foo. We don't want to have to deal with merge conflicts later. Those are
boring and risky. So let's make these changes on top of the changes we made in
step 2:

git append 3-rename-bar

We end up with this branch hierarchy:

master\
1-fix-typos\
2-rename-foo\
3-rename-bar

The corresponding vanilla Git command is

git checkout -b 3-rename-bar

fixing more typos

While renaming bar, we stumbled on a few more typos. Let's add them to the first
branch.

git checkout 1-fix-typos

# make the changes and commit them here

git checkout 3-rename-bar

Back on branch 3-rename-bar, the freshly fixed typos are visible again because
the commit to fix them only exists in branch 1-fix-typos right now. Luckily, Git
Town can propagate these changes through all other branches automatically:

git sync

The corresponding vanilla Git commands are:

git checkout -b 2-rename-foo git merge 1-fix-typos git push git checkout
3-rename-bar git merge 2-rename-foo git push

generalize the infrastructure

Okay, where were we? Right! With things properly named it is now easier to make
sense of larger changes. We cut branch 4-generalize-infrastructure and perform
the refactor in it. It has to be a child branch of 3-rename-bar, since the
improved variable naming done before will make the larger changes we are about
to do now more intuitive.

git append 4-generalize-infrastructure

Again, in vanilla Git:

git checkout -b 4-generalize-infrastructure

Lots of coding and committing into this branch to generalize things. Since
that’s all we do here and nothing else, it’s pretty straightforward to get
through it, though. Off goes the code review for those changes. Shipping the
typo fixes

In the meantime, we got the approval for the typo fixes in step 1. Let’s ship
them!

git ship 1-fix-typos

The vanilla Git commands:

git stash -u # move open changes out of the way git checkout master # update
master so that we ship our changes # on top of the most current changes git pull
git checkout 1-fix-typos # make sure the local machine # has all the changes
made in the # 1-fix-typos branch git pull git merge master # resolve any merge
conflicts # between our feature and the latest master now, # on the feature
branch git checkout master git merge — squash 1-fix-typos # use a squash merge #
to remove all temporary commits # on the branch git push # make our shipped
feature visible to # all other developers git branch -d 1-fix-typo # delete the
shipped branch # from the local machine git push origin :1-fix-typo # delete the
shipped branch # from the remote repository git checkout
4-generalize-infrastructure # return to the branch # we were working on git
stash pop # restore open changes we were working on

With branch 1-fix-typos shipped, our branch hierarchy now looks like this:

master\
2-rename-foo\
3-rename-bar\
4-generalize-infrastructure

synchronizing our work with the rest of the world

We have been at it for a while. Other developers on the team have shipped things
too, and technically the branch 2-rename-foo still points to the previous commit
on master. We don't want our branches to deviate too much from what’s happening
on the master branch since that can lead to more severe merge conflicts later.
Let's get everything in sync!

git sync

The corresponding vanilla Git commands are:

git stash -u # move open changes out of the way git checkout master git pull git
checkout 2-rename-foo git merge master git push git checkout 3-rename-bar git
merge 2-rename-foo git push git checkout 4-generalize-infrastructure git merge
3-rename-bar git push git stash pop # restore what we were working on

Because we used git append to create the new branches, Git Town knows which
branch is a child of which other branch, and can do the merges in the right
order. building the new feature

Back to business. With the new generalized code architecture in place, we can
now add the new feature in a clean way. To build the new feature on top of the
new infrastructure:

git append 5-add-feature

Let’s stop here. Hopefully, it is clear how Git Town allows to work in several
Git branches in parallel. Let’s review:

    each change happens in its own feature branch
    git append creates a new feature branch on top of your existing work
    git sync keeps all feature branches in sync with the rest of the world - do this several times a day
    git ship ships a feature branch

advantages

Working this way has a number of important advantages:

    Focused changes are easier and faster to create: if you change just one thing, you can do it quickly, make sure it makes sense, and move on to the next issue in another branch. No more getting stuck unsure which of the many changes you did in the last 10 minutes broke the build, and no need to fire up the debugger to resolve this mess.
    They are easier and faster to review: The pull request can have a simple description to summarize it. Reviewers can easily wrap their heads around what changes they are looking at, and make sure they are correct and complete. This is also true if you write code just by yourself.
    Branches containing focused changes cause less merge conflicts than branches with many changes in them. This gives Git more opportunity to resolve merge issues automatically.
    In case you have to resolve merge conflicts manually, they are also easier and safer to resolve because the changes in each branch are more obvious.
    Others can start reviewing parts of your changes sooner because you start submitting pull requests earlier.

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
