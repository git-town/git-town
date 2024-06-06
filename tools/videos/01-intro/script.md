# Screencast 1: Introduction to Git Town

## Part 1: summarize and take-aways

Git Town is a free and open-source application
that adds a few missing commands to Git.

These new Git commands allow you to manage Git branches
much more efficiently than is possible with the standard Git commands
and more accurately than when doing it manually.

Git Town reduces and often completely eliminates merge conflicts.
It can also save you from losing work
after accidentally running a wrong Git command.

While executing Git commands for you,
Git Town stays true to Git's nature
as a flexible and powerful tool
that doesn't force you into one particular way of using it.

Git Town is useful when using it just by yourself
And it really shines in collaborative scenarios
when you write code together with other developers.

Let's see it in action!

## Part 2: creating feature branches

Here is a software project that I work on.
Let's start hacking on a new feature!

(( run 01.sh ))

```
git hack my-feature
```

Let's see what Git Town's "hack" command did here.

(( highlight "git checkout -b my-feature" ))
Most people would simply run "git checkout -b branch-name" to create a new branch.
So does Git Town.

(( highlight "git fetch --prune" ))
Before it does that, it pulls updates from the central code repository.

(( highlight "github.com" in the Git output ))
In this case, that's my Github repo.

(( highlight the SHAs of the downloaded commits on main ))
Other developers have added things to the main branch since I last updated my local Git clone.

As you can already tell by now,
we want to cut our new feature branch from the up-to-date main branch
not the outdated one we had locally.

(( highlight "git rebase origin/main ))
Git Town does that here.

(( highlight "git checkout -b my-feature main ))
It then cuts the feature branch off the now up-to-date main branch.

## Part 3: synchronizing feature branches

``` 
prep the codebase:
- sync the feature branch
- create files "other_file_1" and "other_file_2" on the main branch
- create file "my_file_2" on the tracking branch of the feature branch
- create a local file my_file_1 on the feature branch
- delete the shell command history to avoid displaying unrelated commands: rm ~/.local/share/fish/fish_history
```

Let's start building the feature.

(( type "I am implementing my feature" ))

(( screen: "2 hours later" ))
Okay, I have been hacking on this feature for some time now.
I think I better sync my work with changes that were made in the meantime.
If there have been larger changes,
I better get them them into my feature branch before I move any further.
Otherwise I run the risk of modifying something that was also changed on the main branch
and then I have a merge conflict that will be frustrating and error-prone to resolve.

(( run "git sync" ))

"git sync" also ran a number of Git commands for us.
Let's take a quick look what exactly it did.

(( highlight "git fetch --prune" ))
As always, it starts by pulling the latest updates
from the remote repository down to our local repo.

(( highlight "git stash -u" ))
Because we ran "git sync" in the middle of ongoing work,
it stashes our uncommitted changes away
so that they don't get in the way of what happens next.

(( highlight "git checkout main" ))
Now we switch to the main branch

(( highlight "Your branch is behind origin/main by 1 commit" ))
and because the main branch on this machine doesn't have the latest commits from our coworkers

(( highlight "git rebase origin/main" ))
it pulls these updates down from the remote main branch into our local one.

(( highlight "git checkout my-feature" ))
With the main branch up to date now
we can go back to our feature branch.

(( highlight "Your branch is behind origin/my-feature by 1 commit" ))
Some commits on this branch are not on this machine.
I committed them on my other computer and almost forgot about them.

(( highlight "git merge --no-edit origi/kg-feature-1" ))
But Git Town doesn't and pulls these commits down on my machine as well.

(( highlight "git merge --no-edit main" ))
Next we merge our now up-to-date main branch into our now up-to-date feature branch.
If you prefer to rebase your feature branches, you can configure this.

(( highlight "git push" ))
Git Town pushes the new commits on our feature branch to the remote
so that we have a backup in case something happens to our local files
and in case a teammate also uses our feature branch.

(( highlight "git stash pop" ))
Finally, Git Town restores the stashed away uncommitted changes back into our workspace.
This leaves us exactly where we started the "git sync"
except that our feature branch is now fully in sync with the rest of the world.

Similar to the situation right after "git hack",
this means any changes we do from now on
won't conflict with existing code
because our workspace contains all the latest commits now.

If we go back to our editor and refresh the file view

(( refresh editor, highlight the new files ))
we see all the new files that were pulled in by the last "git sync".
They are in our feature branch now.

## Part 4: submitting a proposal

The feature is done, let's submit a pull request.
Since some people call this "merge request", Git Town uses a neutral name.

(( run "git propose" ))

When running "git propose",
Git Town opens the form to submit a pull request in my browser
and populates it with the data it knows.

I can fill in the missing information.

(( fill in the request title and body, select a reviewer ))

(( click the "submit" button )) And off it goes.

## Part 5: cleanup after shipping

The pull request gets reviewed, approved, and merged the normal way.

(( click the "merge" button of the PR ))

Merging it deletes the remote part of my feature branch.
I still need to delete my local copy.
Let's do that now.
I am already working on the next feature.

(( run "git sync --all" ))
When I sync the next time,

(( highlight "removed branch my-feature" ))
Git Town removes the shipped branch from my workspace.
It syncs it locally again to verify
that my local copy of the feature branch does not contain any unshipped changes.
