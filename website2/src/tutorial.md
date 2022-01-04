# Tutorial

Let's assume you are part of a development team with Bob and Carol. Your
repository is hosted on [GitHub](https://github.com) and you follow
[GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow). We
use `main` as the name for the main development branch. Other commonly used
names for this branch are `master` or `development`.

### Starting a new feature

You are in the middle of the sprint and have just finished a feature. You take
the next ticket from the backlog. Let's say it is called "resetting passwords".
Since we are developing in feature branches, you now need to

    update your `main` branch to the latest version on GitHub
    cut a new feature branch from the `main` branch

Running git town hack reset-passwords achieves all this in a single command.
This gives you the best possible start for building the new feature, avoiding
unnecessary merge conflicts because you build on top of the latest version of
the code base.

### Synchronizing the branch

After coding for a while you hear that Bob shipped a number of important bug
fixes and Carol shipped some UI design updates. Both changes affect your work,
so you want them in your feature branch before you continue making more changes
on your end. On a high level, you need to

- pull updates for the `main` branch (to get Bob and Carol's changes)
- merge the `main` branch into your reset-passwords branch
- push your updated feature branch to the repo, so that others who work on it
  also get these updates

You will need to checkout various branches to accomplish this. If you are
currently coding, you likely have uncommitted changes. You should
[stash](https://git-scm.com/docs/git-stash) them away so that they don't get in
the way when changing branches. Altogether this simple operation that you should
run multiple times per day for each of your feature branches requires between 5
and 7 individual Git commands. That's a lot of distraction and typing while
coding.

`git town sync` runs this whole process with a single command and brings you
back to exactly where you started.

With Bob's bug fixes and the new UI from Carol available in your branch, any
more modification you make go right on top of their work. More merge conflicts
have been avoided!

### Creating a pull request

Once your feature is ready for review, it is time to open a pull request on
GitHub. You need to fire up a browser, go to GitHub, navigate to your repository
and finally create a new pull request using the web UI.

`git town new-pull-request` lets you jump straight from the terminal to filling
in the details of your pull request in your browser. Your current branch is
prepopulated, so all you need to do is fill out the title and description, tag
the reviewers, and submit.

### Shipping the feature

When your pull request gets the approval to be merged, you want to ship it. To
do this safely, i.e. without breaking the `main` branch, you want to

- make sure there are no open changes (i.e. all changes are properly committed)
- pull updates from your remote feature branch (to make sure you ship everything
  that is in that branch)
- pull updates for the `main` branch (to make sure you ship on top of the latest
  version of `main`)
- merge the `main` branch into the `password-reset` branch (to make sure your
  branch doesn't create conflicts with`main`, and to give you a chance to
  resolve any issues before merging into`main`)
- squash-merge the `password-reset` branch into the `main` branch (this makes it
  look like a single, clean commit, without the convoluted merge history and the
  many intermediate commits on your branch)
- push the updated `main` branch to the repository (so that your changes are
  available to Bob and Alice in return)
- remove the `password-reset` branch from your local machine and the repository

This requires between 10-15 individual Git commands. Git Town provides a single,
convenient command for this as well: `git town ship` while being on the
`password-reset` branch.

After running this, your feature is now safely merged as a single additional
commit on the `main` branch. The old feature branch is removed from your machine
and from GitHub.
