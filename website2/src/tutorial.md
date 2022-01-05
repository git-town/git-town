# Tutorial

<br>
<p align="center">
  <a href="https://youtu.be/4QDgQajqxRw" target="_blank">
    <img src="video.jpg" width="517" height="290" alt="screencast">
  </a>
</p>
<br>

Let's assume you are part of a development team with Bob and Carol. Your
repository is hosted on [GitHub](https://github.com) and you follow
[GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow). We
use `main` as the name for the main development branch. Other commonly used
names for the main development branch are `master` or `development`.

### Starting a new feature

You are in the middle of the sprint and have just finished a feature. You take
the next ticket from the backlog. Let's say it is called "resetting passwords".
Since we are developing in feature branches, you now need to

- update your `main` branch to the latest version on GitHub
- cut a new feature branch from the `main` branch

Running `git town hack reset-passwords` achieves all this in a single command.
This gives you the best possible start for building the new feature, avoiding
unnecessary merge conflicts because you build on top of the latest version of
the code base.

### Synchronizing the branch

After coding for a while you hear that Bob shipped a number of important bug
fixes and Carol shipped some UI design updates. Both changes potentially affect
your work, so you want them in your feature branch before you continue making
more changes on your end. On a high level, you need to

- pull updates for the `main` branch (to get Bob and Carol's changes)
- merge the `main` branch into your `reset-passwords` branch
- push your updated feature branch to the repo so that others who work on it
  also get these updates

You will need to checkout various branches to accomplish this. If you are
currently coding, you likely have uncommitted changes. You need to
[stash](https://git-scm.com/docs/git-stash) them away so that they don't get in
the way when changing branches. Altogether this simple operation requires
between 5 and 7 individual Git commands. And you should run this multiple times
per day, for each of your feature branches, to avoid larger merge issues. That's
a lot of distraction.

`git town sync` runs this whole process with a single command and brings you
back to exactly where you started. `git town sync --all` syncs all branches on
your machine at once.

### Creating a pull request

Once your feature is ready for review, it is time to open a pull request on
GitHub. You need to do a final sync of your branch, fire up a browser, go to
GitHub, navigate to your repository and create a new pull request using the web
UI.

`git town new-pull-request` lets you jump straight from the terminal to filling
in the details of your pull request in your browser. Your current branch is
prepopulated, so all you need to do is fill out the title and description, tag
the reviewers, and submit.

### Shipping the feature

When your pull request gets the approval to be merged, you "ship" it by merging
it into the main code branch so that your changes can start their journey into
production. To do this safely, i.e. without breaking the build on the `main`
branch, you want to

- make sure you didn't forget to commit any changes in your workspace
- pull updates from your remote feature branch to make sure you ship everything
  that is in that branch
- pull updates for the `main` branch to make sure you ship on top of the latest
  version of `main`
- merge the `main` branch into your feature branch to make sure your branch
  doesn't create conflicts with `main` and to give you a chance to resolve any
  issues before merging into `main`
- squash-merge the `password-reset` branch into the `main` branch. This makes it
  look like a single, clean commit, without the convoluted merge history and the
  many intermediate commits on your branch
- push the updated `main` branch to the repository so that your changes are
  available to your co-workers
- remove the `password-reset` branch from your local machine and the repository

This requires between 10-15 individual Git commands. Git Town does all this for
you when running `git town ship` while being on the `password-reset` branch.

After running this, your feature is now safely merged as a single additional
commit on the `main` branch. The old feature branch is removed from your machine
and from GitHub.
