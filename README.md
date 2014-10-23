## _Git Automation for Agile Development Teams_
<a href="https://travis-ci.org/Originate/git-town" alt="Build Status" target="_blank"><img src="https://travis-ci.org/Originate/git-town.svg?branch=master"></a>

Git Town provides a number of additional Git commands that
automate the typical high-level operations in
<a href="http://scottchacon.com/2011/08/31/github-flow.html" target="_blank">GitHub flow</a>
and others.

It is designed for workflows that have a main branch
(typically "development" or "master")
from which feature branches are cut and into which they are merged,
and assumes you use a central code repository like
<a href="http://github.com" target="_blank">GitHub</a> or
<a href="https://bitbucket.org" target="_blank">BitBucket</a>.

* create a new feature branch with <a href="#git-hack">git hack</a>
* keep your feature branch in sync with the rest of the world using <a href="#git-sync">git sync</a>
* when done with a feature, merge it into the main branch with <a href="#git-ship">git ship</a>
* keep a forked repository up to date with its upstream using <a href="#git-sync-fork">git sync-fork</a>
* extract existing commits into their own feature branches using <a href="#git-extract">git extract</a>
* delete merged branches in local and remote repository with <a href="#git-prune-branches">git prune-branches</a>.

Git Town automatically prunes no longer existing remote branches
from your branch list.

Hint: the examples below assume "master" as the main branch name
(this is <a href="#configuration">configurable</a>),
and "feature" as the feature branch name.


## git hack

_Cuts a new feature branch off the main branch._

Scenario:
While working on something you realize "Hey, this should be in its own branch."
No problem, just run `git hack foo`,
and you get an up to date feature branch "foo"
with all uncommitted changes copied into it.

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>checkout the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull main branch updates from the repo</td>
    <td>git fetch<br>git rebase origin/master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>cut the new feature branch</td>
    <td>git checkout -b feature master</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>restore the stashed changes</td>
    <td>git stash pop</td>
  </tr>
</table>



## git sync

_Syncronizes the current feature branch with the rest of the world._

This works even when you are in the middle of coding,
with uncommitted changes in your workspace.
You can call this command safely at any time, many times during the day.

* run the command: `git sync`
* abort the command when there are conflicts: `git sync --abort`
* finish the sync after you have fixed the conflicts: `git sync --continue`
* pushes tags when run the main or a non-feature branch
* skips unnecessary pushes

_on a feature branch:_

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>checkout the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull main branch updates from the repo</td>
    <td>git fetch<br>git rebase origin/master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>checkout the feature branch</td>
    <td>git checkout feature</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>pull feature branch updates from the repo</td>
    <td>git merge origin/feature</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>update the feature branch with the latest changes from main</td>
    <td>git merge master</td>
  </tr>
  <tr>
    <td>7a.</td>
    <td>push the feature branch (if we don't have a remote branch yet)</td>
    <td>git push -u origin feature</td>
  </tr>
  <tr></tr>
  <tr>
    <td>7b.</td>
    <td>push the feature branch (with existing remote branch)</td>
    <td>git push</td>
  </tr>
  <tr>
    <td>8.</td>
    <td>restore the stashed changes</td>
    <td>git stash pop</td>
  </tr>
</table>

_on the main or a non-feature branch:_

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>pull branch updates from the repo</td>
    <td>git fetch<br>git rebase origin/[branch name]</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>push the branch</td>
    <td>git push</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>push tags</td>
    <td>git push --tags</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>restore the stashed changes</td>
    <td>git stash pop</td>
  </tr>
</table>


## git extract

_Extracts commits from a feature branch into another._

Scenario:
While working on a bigger feature you want to extract certain changes
(like refactorings) into their own feature branches, so that they can be
reviewed separately/faster than the rest of the feature.

More background around
<a href="http://blog.originate.com/blog/2014/04/19/refactoring_git_branches" target="_blank">Git branch refactoring</a>.

* run the command: `git extract [new branch name]`<br>
* abort the command when there are conflicts: `git extract --abort`
* you need <a href="http://en.wikipedia.org/wiki/Ncurses" target="_blank">Ncurses</a> for this

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>checkout the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull main branch updates from the repo</td>
    <td>git fetch<br>git rebase origin/master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>user picks the commits to extract</td>
    <td>(nice GUI tool)</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>cut a new feature branch off the main branch</td>
    <td>git checkout -b new_feature master</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>copy the chosen commits to the feature branch</td>
    <td>git cherry-pick [SHA1 of the commits]
  </tr>
  <tr>
    <td>7.</td>
    <td>restore the stashed away changes</td>
    <td>git stash pop</td>
  </tr>
</table>


## git ship

_Ships a finished feature._

Call this from the feature branch that you want to ship.

* run the command: `git ship`
* run the command passing in the squashed commit message: `git ship -m [commit message]`
* abort the command when there are conflicts: `git ship --abort`
* verifies that we are shipping a feature branch, that the feature branch
  has shippable commits, and that there are no uncommitted changes.

<table>
  <tr>
    <td>1.</td>
    <td>ensure there are no uncommitted changes</td>
    <td>git status</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>checkout the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull main branch updates from the repo</td>
    <td>git fetch<br>git rebase origin/master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>checkout the feature branch</td>
    <td>git checkout feature</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>pull updates for the feature branch</td>
    <td>git merge origin/feature</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>checkout the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>7.</td>
    <td>merge the feature branch into the main branch</td>
    <td>git merge --squash feature</td>
  </tr>
  <tr>
    <td>8.</td>
    <td>push the updated main branch</td>
    <td>git push</td>
  </tr>
  <tr>
    <td>9.</td>
    <td>delete the feature branch from the developer machine</td>
    <td>git branch -d feature</td>
  </tr>
  <tr>
    <td>10.</td>
    <td>delete the feature branch from the repo</td>
    <td>git push origin :feature
  </tr>
</table>


## git sync-fork

_Syncs the main branch with the upstream repository._

Call this to bring the main branch up to date with the main branch of the remote `upstream`.

If your respository is a fork on GitHub, `upstream` will be automatically set on first use.

* run the command: `git sync-fork`
* skips unnecessary pushes

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>checkout the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull upstream updates for the main branch</td>
    <td>git fetch upstream<br/>git rebase upstream/master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>push the main branch</td>
    <td>git push</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>checkout the branch you started on</td>
    <td>git checkout [initial branch]</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>restore the stashed away changes</td>
    <td>git stash pop</td>
  </tr>
</table>


## git prune-branches

_Delete merged branches in local and remote repository._

* run the command: `git prune-branches`
* if the current branch is merged, moves to the main branch

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>pull updates from the repo</td>
    <td>git fetch</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>delete each merged branch in the remote repository</td>
    <td>git push origin :&lt;branch_name&gt;</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>delete each merged branch in the local repository</td>
    <td>git branch -d &lt;branch_name&gt;</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>restore the stashed away changes</td>
    <td>git stash pop</td>
  </tr>
</table>

## installation

Git Town is 100% bash script, so it runs anywhere where Git and Bash runs:
OS X, Linux, BSD, and even Windows with Cygwin or something similar.
Installation on OS X can be done using <a href="http://brew.sh" target="_blank">Homebrew</a>,
other platforms need to install manually.

<table>
  <tr>
    <th width="300px">
      Using Homebrew
    </th>
    <th width="300px">
      Manually
    </th>
  </tr>
  <tr>
    <td colspan="2">
      Install
    </td>
  </tr>
  <tr>
    <td>
      brew tap Originate/gittown<br>
      brew install git-town
    </td>
    <td>
      clone the repo to your machine<br>
      add the folder to your path
    </td>
  </tr>
  <tr>
    <td colspan="2">
      Update
    </td>
  </tr>
  <tr>
    <td>
      brew update<br>
      brew upgrade git-town
    </td>
    <td>
      cd [directory of your Git Town clone]<br>
      git pull
    </td>
  </tr>
  <tr>
    <td colspan="2">
      Uninstall
    </td>
  </tr>
  <tr>
    <td>
      brew uninstall git-town<br>
      brew untap Originate/gittown
    </td>
    <td>
      remove repo from your machine<br>
      remove folder from path
    </td>
  </tr>
</table>



## configuration

On first use, Git Town will ask for the main branch name and the names of any other non feature branches.
Git Town stores its configuration in the Git configuration of your project.
If these ever need to change, the configuration can be updated using <a href="http://git-scm.com/docs/git-config" target="_blank">git config</a>.


```bash
# Read configuration
git config git-town.main-branch-name
git config git-town.non-feature-branch-names

# Write configuration
git config git-town.main-branch-name master
git config git-town.non-feature-branch-names 'qa, production'
```

## development

* run all tests: `cucumber`
* run a single test: `cucumber -n 'scenario/feature name'` or `cucumber [filename]:[lineno]`

Found a bug or want to contribute a feature?
<a href="https://github.com/Originate/git-town/issues/new" target="_blank">Open an issue</a>
or - even better - get down, go to town, and fire a feature-tested
<a href="https://help.github.com/articles/using-pull-requests" target="_blank">pull request</a>
our way!


## roadmap

The future roadmap is planned using
<a href="https://github.com/Originate/git-town/issues" target="_blank">GitHub issues</a>.
The past roadmap is in the <a href="release-notes.md" target="_blank">release notes</a>.

If you have an idea about a cool feature you would like to see in Git Town,
please <a href="https://github.com/Originate/git-town/issues/new" target="_blank">open a ticket</a>.
Our team will add the <a href="https://github.com/Originate/git-town/labels/idea" target="_blank">idea</a> tag.
Once we reach agreement about this idea, it will be tagged as <a href="https://github.com/Originate/git-town/labels/enhancement" target="_blank">enhancement</a>
or <a href="https://github.com/Originate/git-town/labels/bug" target="_blank">bug</a>.

