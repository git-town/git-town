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
* in a forked repository, keep your main branch up to date with its upstream repository using <a href="#git-sync-fork">git sync-fork</a>
* refactor selected commits from one feature branch into a dedicated one using <a href="#git-extract">git extract</a>.

All Git Town commands automatically clean up (prune)
no longer existing remote branches from your branch list.

Hint: the examples below assume "master" as the main branch name
(this is <a href="#configuration">configurable</a>),
and "feature" as the feature branch name.


## git hack

_Cuts a new feature branch off the main branch._

Scenario:
While working on something you realize "Hey, this should be in its own branch."
No problem, just run `git hack foobar`,
and a feature branch with name "foobar" is created for you,
with all open changes copied over into it.

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>check out the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull updates for the main branch</td>
    <td>git pull</td>
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

_Syncronizes the current feature branch with the rest of the world,
i.e. with its remote branch and the main branch._

This works even when you are right in the middle of coding,
i.e. with uncommitted changes.
You can call this command safely at any time, many times during the day.

* run the command: `git sync`<br>
* abort the command when there are conflicts: `git sync --abort`<br>
* finish the sync after you have fixed the conflicts: `git sync --continue`

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td> git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>pull feature branch updates from the repo</td>
    <td>git fetch<br>git merge origin/feature</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>switch to the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>pull main branch updates from the repo</td>
    <td>git merge origin/master</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>switch to the feature branch</td>
    <td>git checkout feature</td>
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


## git extract

_Extracts commits from a feature branch into a new feature branch._

Scenario:
After finishing a bigger feature you realize that this is actually several
changes in one branch. You want to extract each change into its own feature
branch.

More background around <a href="http://blog.originate.com/blog/2014/04/19/refactoring_git_branches" target="_blank">Git branch refactoring</a>.

* run the command: `git extract [new branch name]`<br>
* abort the command when there are conflicts: `git extract --abort`

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td> git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>switch to the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull updates for the main branch</td>
    <td>git pull</td>
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

<table>
  <tr>
    <td>1.</td>
    <td>ensure there are no uncommitted changes</td>
    <td>git status</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>pull updates for the feature branch</td>
    <td>git fetch<br>git merge origin/feature</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>check out the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>pull updates for the main branch</td>
    <td>git merge origin/master</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>merge the feature branch into the main branch</td>
    <td>git merge --squash feature</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>push the updated main branch</td>
    <td>git push</td>
  </tr>
  <tr>
    <td>7.</td>
    <td>delete the feature branch from the developer machine</td>
    <td>git branch -d feature</td>
  </tr>
  <tr>
    <td>8.</td>
    <td>delete the feature branch from the repo</td>
    <td>git push origin :feature
  </tr>
</table>


## git sync-fork

_Syncs the main branch with the upstream repository._

Call this to bring the main branch up to date with the main branch of the remote `upstream`.

If your respository is a fork on GitHub, `upstream` will be automatically set on first use.

* run the command: `git sync-fork`

<table>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td>git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>check out the main branch</td>
    <td>git checkout master</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull upstream updates for the main branch</td>
    <td>git fetch upstream<br/>git merge upstream/master</td>
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

Git Town asks for the main branch name if one isn't set per repository,
and stores this information in the Git configuration of your project.


## development

* run all tests: `cucumber`
* run a single test: `cucumber -n 'scenario or feature name'` or `cucumber [filename]:[lineno]`


## roadmap

The roadmap is developed using readme-driven development <a href="RDD.md">here</a>.
Also check out the <a href="release-notes.md">release notes</a>.
