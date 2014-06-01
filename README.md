## _Elegant Git workflows for a more civilized age_
<a href="https://travis-ci.org/Originate/git-town" alt="Build Status" target="_blank"><img src="https://travis-ci.org/Originate/git-town.svg?branch=master"></a>

Git Town is a configurable\* and hackable collection of additional Git commands that perform the typical high-level operations which a software developer performs (or should perform) in a collaborative environment.

This includes things like easily creating up-to-date feature branches,
keeping feature branches synchronized with the ongoing development from other developers,
as well as removing feature branches from the developer machine and the central repo after they have been merged with the main code line.

Git Town operates under the following assumptions:

* You have a **main development branch** (typically "development" or "master"). In this documentation we will use "development".
* You follow a per-project strategy that prefers either _rebases_ or _merges_ for updating branches, and _squash merges_ or _normal merges_ for merging feature branches into the main branch\*.
* You use a central code repository like [Github](http://github.com) (called __repo__ from now on).


_Note: This documentation is the driver for readme-driven development of this tool.
The parts marked with an asterisk (\*) are not yet implemented. Feedback for them is welcome!_


## Scripts

Git Town provides the following Git commands.


### git hack

_Cuts a new feature branch off the main development branch. Even if you are right in the middle of something._

Run the command: `git hack [name of feature branch to create]`

<table>
  <tr>
    <th colspan="2">step</th>
    <th>rebase version</th>
    <th>merge version**</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td colspan="2" align="center">git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>check out the main branch</td>
    <td colspan="2" align="center">git checkout development</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull updates for the main branch from the repo</td>
    <td>git pull --rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>cut a new feature branch off the main branch</td>
    <td colspan="2" align="center">git checkout -b [feature branch] development</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>restore the stashed changes</td>
    <td colspan="2" align="center">git stash pop</td>
  </tr>
</table>



### git sync

_Syncronizes the current feature branch with the rest of the world, i.e. with its remote branch and the main branch.
This also works when you are right in the middle of something._

Run the command: `git sync`<br>
Abort the command when there are conflicts: `git sync --abort`

<table>
  <tr>
    <th colspan="2" align="center">step</th>
    <th>rebase version</th>
    <th>merge version*</th>
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
    <td colspan="2" align="center">git checkout development</td>
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
    <td>git rebase development</td>
    <td>git merge development</td>
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


### git sync --all\*

_Synchronizes all branches on the local machine with the rest of the world._

* does a `git sync` on each feature branch that exists on the local machine


### git extract

_Extracts commits from a feature branch into a new feature branch._

More background around <a href="http://blog.originate.com/blog/2014/04/19/refactoring_git_branches" target="_blank">Git branch refactoring</a>.

Run the command: `git extract [new branch name]`<br>
Abort the command when there are conflicts: `git extract --abort`

<table>
  <tr>
    <th colspan="2" align="center">step</th>
    <th>rebase version</th>
    <th>merge version*</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>optionally stash uncommitted changes away</td>
    <td colspan="2" align="center"> git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>switch to the main branch</td>
    <td colspan="2" align="center">git checkout development</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>pull the latest updates for the main branch from the repo</td>
    <td>git pull --rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>lets the user pick the commits to extract</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>cut a new feature branch off the main branch</td>
    <td colspan="2" align="center">git checkout -b [feature] development</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>cherry-pick the selected commits into the new branch</td>
    <td colspan="2" align="center">git cherry-pick [SHA1 of the commits]
  </tr>
  <tr>
    <td>7.</td>
    <td>restore the stashed away changes</td>
    <td colspan="2" align="center">git stash pop</td>
  </tr>
</table>


### git ship

_Ships a finished feature._

When on the feature branch to ship, run the command: `git ship`<br>
Abort the command when there are conflicts: `git ship --abort`

<table>
  <tr>
    <th colspan="2" align="center">step</th>
    <th>rebase version</th>
    <th>merge version*</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>ensure there are no uncommitted changes in the workspace</td>
    <td colspan="2" align="center">git status</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>pull updates for the feature branch from the repo</td>
    <td>git pull<br>--rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>check out the main branch</td>
    <td colspan="2" align="center">git checkout development</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>pull updates for the main branch from the repo</td>
    <td>git pull<br>--rebase</td>
    <td>git pull</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>merge the feature branch into the main branch</td>
    <td colspan="2" align="center">git merge --squash [feature branch]
  </tr>
  <tr>
    <td>6.</td>
    <td>push the new commit of the main branch to the source repo</td>
    <td colspan="2" align="center">git push</td>
  </tr>
  <tr>
    <td>7.</td>
    <td>delete the feature branch from the local machine</td>
    <td colspan="2" align="center">git branch -d [feature branch]</td>
  </tr>
  <tr>
    <td>8.</td>
    <td>delete the feature branch from the remote repo</td>
    <td colspan="2" align="center">git push origin :[feature branch]
  </tr>
</table>


### git kill
Safely deletes a git branch.

* never deletes the main branch
* if the branch has unmerged commits, asks the user for confirmation\*
* deletes the given branch from the local machine as well as the repo


## Configuration\*

Each Git Town command comes in a _rebase_ and a _merge_ version.
Which option is used can be configured through the following options in ".git-town-rc"


### Pull Strategy\*
The pull strategy defines which command Git Town uses when updating a branch from its remote tracking branch.

* __rebase__: always do a `git pull --rebase`, let the user figure out merge conflicts
* __automatic__: try `git pull --rebase` first. When merge conflicts occur, abort and do a `git pull`
* __merge__: always do `git pull`

The default value for this setting is _automatic_.


### Update Strategy\*
The update strategy defines which command Git Town uses when updating a feature branch with updates from the main branch.

* __rebase__: do a `git rebase [main branch]`, let the user resolve eventual merge conflicts. This setting should only be used if your feature branches are always private. It results in `git push --force` when pushing the updates back to the remote tracking branch.
* __automatic__: try a `git rebase [main branch]` first. If merge conflicts occur abort, abort and try a `git merge [branch name]`
* __merge__: always do `git merge [main branch]`

The default value for this setting is _automatic_.

### Merge Strategy\*
The merge strategy defines which command Git Town uses when merging feature branches into the main branch.

* __squash__: always do `git merge --squash`
* __no-ff__ : always do `git merge --no-ff`
* __normal__: always do `git merge`


## Installation


<table>
  <tr>
    <th>
      Using Homebrew
    </th>
    <th>
      Manually
    </th>
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
</table>


#### Updating

<table>
  <tr>
    <th>
      Using Homebrew
    </th>
    <th>
      Manually
    </th>
  </tr>
  <tr>
    <td>
      brew update<br>
      brew upgrade git-town
    </td>
    <td>
      git pull
    </td>
  </tr>
</table>


#### Uninstalling

<table>
  <tr>
    <th>
      Using Homebrew
    </th>
    <th>
      Manually
    </th>
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



## Develop your own scripts

* run all tests: `spec/run`
* run a single test: `spec/run [test filename]`
* this script clones https://github.com/Originate/git_town_specs into your `/tmp` directory


Some background on the code structure:
* Due to limitations of Bash Script, the functions take normal arguments, and return their result as global variables.
* Each function does the thing it says in a robust way. The "pull_feature_branch" function for example switches to the current feature branch, and then pulls it.

