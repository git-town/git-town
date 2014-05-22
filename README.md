## _Elegant Git workflows for a more civilized age_
[![Build Status](https://travis-ci.org/Originate/git-town.svg?branch=master)](https://travis-ci.org/Originate/git-town)

Git Town is an easily hackable collection of additional Git commands that make massively parallelized collaborative software development easy, safe, and fun.

* You have a **main development branch** (typically "development" or "master").
* You follow a strategy that prefers either _rebases_ or _merges_.
* Your team decided to always use _squash merges_ - or not.
* You use a central code repository like [Github](http://github.com) (called __repo__ from now on).


_Note: This documentation is the driver for readme-driven development of this tool.
The parts marked with an asterisk are not yet implemented._


## Configuration\*

Each Git Town command comes in a _rebase_ and a _merge_ version.
Which option is used can be configured through the following options in ".git-town-rc"
* __rebase=always__: never fall back to merging, let the user figure out merge conflict or abort
* __rebase=automatic__: try rebasing first, when merge conflicts occur abort, undo, and retry in the merge version
* __rebase=never__: always do merges

The default value for the global setting of _rebase_ is _always_.
You can override it globally, and per branch:
* always_merge_branches=(master development)

You can also configure how Git Town should perform merges for you:
* __merge=squash__: do _squash merges_
* __merge=no-ff__ : do _-no-ff merges_
* __merge=normal__: do _normal merges_


## Scripts

Git Town provides the following extra Git commands.


### git hack

_Cuts a new feature branch off the main development branch.<br>Even when you are right in the middle of something._

<table>
  <tr>
    <th colspan="2">step</th>
    <th>rebase command</th>
    <th>merge command**</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td colspan="2" align="center">git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>check out the main development branch</td>
    <td colspan="2" align="center">git checkout [main branch]</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull the latest updates for the main branch from the repo</td>
    <td>git pull --rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>cut a new feature branch off the main branch</td>
    <td colspan="2" align="center">git checkout -b [feature branch] [main]</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>restore the stashed changes</td>
    <td colspan="2" align="center">git stash pop</td>
  </tr>
</table>

* run the command: `git hack [name of feature branch to create]`
* abort and undo the command if there are conflicts\*: `git hack --abort`
* abort and retry the command in merge mode\*: `git hack --retry --merge`


### git sync

_Syncronizes the current feature branch with the rest of the world.<br>Even when you are right in the middle of something._

<table>
  <tr>
    <th colspan="2" align="center">step</th>
    <th>rebase command</th>
    <th>merge command*</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td colspan="2" align="center"> git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>pull the latest updates for the feature branch</td>
    <td>git pull<br>--rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>switch to the main branch</td>
    <td colspan="2" align="center">git checkout [main]</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>pull the latest updates for the main branch from the repo</td>
    <td>git pull<br>--rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>switch to the feature branch</td>
    <td colspan="2" align="center">git checkout [feature]</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>update the feature branch with the latest updates from the main branch</td>
    <td>git rebase [main]</td>
    <td>git merge [main]</td>
  </tr>
  <tr>
    <td>7a.</td>
    <td>if there is no remote branch, push the updated feature branch to the code repo and set up branch tracking</td>
    <td colspan="2" align="center">git push -u origin [feature]</td>
  </tr>
  <tr>
    <td>7b.</td>
    <td>if there is a remote branch, push the updated feature branch to the code repo</td>
    <td>git push<br>--force</td>
    <td>git push</td>
  </tr>
  <tr>
    <td>8.</td>
    <td>restore the stashed away changes</td>
    <td colspan="2" align="center">git stash pop</td>
  </tr>
</table>

* run the command: `git sync`
* abort the command when there are conflicts\*: `git sync --abort`


### git sync --all\*

_Synchronizes all branches on the local machine with the rest of the world._

* does a `git sync` on each feature branch that exists on the local machine


### git extract

_Extracts commits from a feature branch into another feature branch._

More background around <a href="http://blog.originate.com/blog/2014/04/19/refactoring_git_branches" target="_blank">Git branch refactoring</a>.

<table>
  <tr>
    <th>step</th>
    <th>rebase command</th>
    <th>merge command*</th>
  </tr>
  <tr>
    <td>1. optionally stash uncommitted changes away</td>
    <td colspan="2" align="center"> git stash</td>
  </tr>
  <tr>
    <td>2. switch to the main branch</td>
    <td colspan="2" align="center">git checkout main</td>
  </tr>
  <tr>
    <td>3. pull the latest updates for the main branch from the repo</td>
    <td>git pull --rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>lets the user pick the commits to extract</td>
  </tr>
  <tr>
    <td>cut a new feature branch off the main branch</td>
    <td colspan="2" align="center">git checkout -b [feature] [main]</td>
  </tr>
  <tr>
    <td>cherry-pick the selected commits into the new branch</td>
    <td colspan="2" align="center">git cherry-pick [SHA1 of the commits]
  </tr>
  <tr>
    <td>8. restore the stashed away changes</td>
    <td colspan="2" align="center">git stash pop</td>
  </tr>
</table>


### git ship

_Ships a finished feature._

<table>
  <tr>
    <th>step</th>
    <th>rebase command</th>
    <th>merge command*</th>
  </tr>
  <tr>
    <td>ensure there are no uncommitted changes in the workspace</td>
    <td colspan="2" align="center">git status</td>
  </tr>
  <tr>
    <td>pull the latest updates for the main branch from the repo</td>
    <td>git pull --rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>squash-merge the current feature branch into the main branch</td>
    <td colspan="2" align="center">git merge --squash [feature branch]
  </tr>
  <tr>
    <td>push the new commit of the main branch to the source repo</td>
    <td colspan="2" align="center">git push</td>
  </tr>
  <tr>
    <td>delete the feature branch from the local machine</td>
    <td colspan="2" align="center">git br -d [feature branch]</td>
  </tr>
  <tr>
    <td>delete the feature branch from the remote repo</td>
    <td colspan="2" align="center">git push origin :[feature branch]
  </tr>
</table>


### git undo\*

Undoes the last Git Town operation.

* git hack: remove the new feature branch and return to the previous feature branch
* git extract: delete the new feature branch and return to the previous feature branch


## Installation

### Using Homebrew\*

`brew install git-town`


### Manually

* clone the repo to your machine
* add the folder to your path
* get busy


## Develop your own scripts

* check out the existing scripts like https://github.com/Originate/git-town/blob/master/git-hack
* check out the available helpers at https://github.com/Originate/git-town/tree/master/helpers
* add more helpers
* add more scripts
* share useful stuff back as a pull request

