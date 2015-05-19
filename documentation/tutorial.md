# Git Town Tutorial

Let's assume you are part of a development team with Bob and Carol.
Your repository is hosted on GitHub, and you follow [GitHub Flow](https://guides.github.com/introduction/flow/index.html).


## Starting a new feature

You are in the middle of the sprint and have just finished a feature.
You take the next ticket from the backlog.
Let's say it is called "resetting passwords".
Since we are developing in feature branches, you now need to

* update your master branch to the latest version on GitHub
* cut a new feature branch from the master branch

Running `git hack reset-passwords` achieves all this in a single command.
This gives you the best possible start for building the new feature,
on top of the latest version of the code base.


## Synchronizing the branch

After coding for a while you hear that Bob shipped a number of important bug fixes,
and that Carol shipped some UI design updates.
Both changes affect your work, so you want them in your feature branch before you continue.
On a high level, you need to

* pull updates for the `master` branch (to get Bob and Carol's changes)
* merge the master branch into your `reset-passwords` branch
* push your updated feature branch to the repo, so that others who work on it also get these updates

You will need to move between branches to do this,
which means you also need to stash away any currently open changes in your repo temporarily.
Altogether this simple operation requires between 5 and 7 individual Git commands.
That's a lot of typing.
And this can (should) happen several times per day, for each of your feature branches!

`git sync` runs this whole process with a single command and brings you back to exactly where you started.

With Bob's bug fixes and the new UI from Carol available in your branch,
any more modification you make go right on top of their work.


## Creating a pull request

Once your feature is ready for review, it's time to open a pull request on GitHub.
You fire up a browser, go to GitHub, navigate to your repository and finally end up on a new pull request form.

`git new-pull-request` lets you jump straight from the terminal to filling in the details of your pull request in your browser.
Your current branch is already selected,
so all you need to do is fill out the title and description,
tag the reviewers, and submit.


## Shipping the feature

When your pull request gets the approval to be merged,
you want to ship it.
To do this safely, i.e. without breaking the master branch, you want to

* make sure there are no open changes (i.e. all changes are properly committed)
* pull updates from your remote feature branch (to make sure you ship everything that is in that branch)
* pull updates for the master branch (to make sure you ship on top of the latest version of master)
* merge the _master_ branch into the _password-reset_ branch
  (to make sure your branch doesn't create conflicts with _master_,
  and to give you a chance to resolve any issues before merging into _master_)
* squash-merge the _password-reset_ branch into the _master_ branch (this makes it look like a single, clean commit, without the convoluted merge history and the many intermediate commits on your branch)
* push the updated _master_ branch to the repository (so that your changes are available to Bob and Alice in return)
* remove the _password-reset_ branch from your local machine and the repository

This requires from 7-9 individual Git commands.
Git Town provides a single, convenient command for this as well:
`git ship reset-passwords`.

After running this, your feature is now safely merged as a single additional commit on the _master_ branch,
then the old feature branch is cleaned up everywhere,
and you and the repository are ready for the next feature.
