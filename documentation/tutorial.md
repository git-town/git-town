# Git Town Tutorial

Let's assume your name is Jessie, and you are part of a development team with Bob and Carol. Your repository is hosted on GitHub, and you follow [GitHub Flow](https://guides.github.com/introduction/flow/index.html).


## Starting a new feature

You are in the middle of the sprint, and have just finished a feature. After a short stretch you take the next ticket from the backlog. It's called "resetting passwords".

You run `git hack password-reset`. This command

* checks out the master branch
* pulls updates for the master branch, i.e. the things that Bob and Carol have shipped while you worked on your last feature
* cuts a new feature branch called "password-reset" from your now up-to-date master branch
* checks out that new branch

This gives you the best possible start for the new password reset feature.


## Synchronizing the branch

After coding for a while, you overhear that Bob shipped a number of important bug fixes, and that Carol shipped some UI updates. Both affect your work, so you want them in your branch.

In order to achieve that, you run `git sync`. This command

* stashes away your currently open changes
* checks out the master branch
* pulls the updates on that branch (Bob's bug fixes and Carol's UI updates)
* checks out your feature branch again
* merges the updates from the master branch into your branch
* pushes your updated feature branch to the repository
* restores your open changes by popping the stash

You are now exactly where you were before, but your branch now also contains Bob's bug fixes and the new UI from Carol. Any more changes you make will fit right in. Great team work!


## Creating a pull request

When you are done, you run `git pr`.
It opens your browser with the GitHub page for creating a new pull request.
You fill out the title and description, tag the reviewers, and submit.


## Shipping the feature

After a while, your pull request gets the approval to be merged. You run `git ship password-reset`. This command

* checks out the master branch
* pulls updates from the remote master branch (to make sure you ship on top of the latest version of master)
* checks out your "password-reset" branch
* pulls updates from its remote branch (to make sure you ship everything that is in that branch)
* merges the master branch into the password-reset branch (to make sure your branch doesn't create conflicts with master, and to give you a chance to resolve them on your branch)
* checks out the master branch again
* squash-merges the password-reset branch into the master branch (this makes it look like a single commit, without the convoluted merge history and the many intermediate commits on your branch)
* pushes the updated master branch to the repository
* deletes the password-reset branch from your local machine and the repository

What all of this achieves is that your feature is safely merged as a single additional commit onto the master branch.
