## _Git Automation for Agile Development Teams_
<a href="https://travis-ci.org/Originate/git-town" alt="Build Status" target="_blank"><img src="https://travis-ci.org/Originate/git-town.svg?branch=master"></a>

* provides a number of additional Git commands
* automates the typical Git operations in <a href="http://scottchacon.com/2011/08/31/github-flow.html" target="_blank">GitHub flow</a> (and others)
* does all the extra updates on each step to keep all branches in sync at all times
* configurable\* and easily hackable

_Note: This documentation is the driver for the readme-driven development of this tool.
The parts marked with an asterisk (\*) are not yet implemented. Feedback is welcome!_


## Scripts

Git Town provides the following Git commands.

_Hint: This documentation uses "master" as the main branch name, and "feature" as the feature branch name._


### git hack

_Cuts a new feature branch off the main branch. Even if you are right in the middle of something._

Run the command: `git hack [name of feature branch to create]`

<table>
  <tr>
    <th colspan="2">step</th>
    <th>rebase version</th>
    <th>merge version*</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td colspan="2" align="center">git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>check out the main branch</td>
    <td colspan="2" align="center">git checkout master</td>
  </tr>
  <tr>
    <td rowspan="3">3.</td>
    <td rowspan="3">pull updates for the main branch</td>
    <td rowspan="2" colspan="2" align="center">git fetch</td>
  </tr>
  <tr></tr>
  <tr>
    <td>git rebase origin/master</td>
    <td>git merge origin/master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>cut the new feature branch</td>
    <td colspan="2" align="center">git checkout -b feature master</td>
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
    <th colspan="2">step</th>
    <th width="28%">rebase version</th>
    <th width="28%">merge version*</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>stash away uncommitted changes</td>
    <td colspan="2" align="center"> git stash</td>
  </tr>
  <tr>
    <td rowspan="3">2.</td>
    <td rowspan="3">pull feature branch updates</td>
    <td rowspan="2" colspan="2" align="center">git fetch</td>
  </tr>
  <tr></tr>
  <tr>
    <td>git rebase origin/feature</td>
    <td>git merge origin/feature</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>switch to the main branch</td>
    <td colspan="2" align="center">git checkout master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>pull main branch updates</td>
    <td>git rebase origin/master</td>
    <td>git merge origin/master</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>switch to the feature branch</td>
    <td colspan="2" align="center">git checkout feature</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>update feature branch</td>
    <td>git rebase master</td>
    <td>git merge master</td>
  </tr>
  <tr>
    <td>7a.</td>
    <td>push the feature branch (if we don't have a remote branch yet)</td>
    <td colspan="2" align="center">git push -u origin feature</td>
  </tr>
  <tr></tr>
  <tr>
    <td>7b.</td>
    <td>push the feature branch<br>(with remote branch)</td>
    <td>git push --force</td>
    <td>git push</td>
  </tr>
  <tr>
    <td>8.</td>
    <td>restore the stashed changes</td>
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
    <th width="28%">rebase version</th>
    <th width="28%">merge version*</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>stash uncommitted changes</td>
    <td colspan="2" align="center"> git stash</td>
  </tr>
  <tr>
    <td>2.</td>
    <td>switch to the main branch</td>
    <td colspan="2" align="center">git checkout master</td>
  </tr>
  <tr>
    <td rowspan="3">3.</td>
    <td rowspan="3">pull updates for the main branch</td>
    <td rowspan="2" colspan="2" align="center">git fetch</td>
  </tr>
  <tr></tr>
  <tr>
    <td>git rebase origin/master</td>
    <td>git merge origin/master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>user picks the commits to extract</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>cut a new feature branch off main</td>
    <td colspan="2" align="center">git checkout -b feature master</td>
  </tr>
  <tr>
    <td>6.</td>
    <td>copy the chosen commits over</td>
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
    <th width="29%">rebase version</th>
    <th width="28%">merge version*</th>
  </tr>
  <tr>
    <td>1.</td>
    <td>ensure no uncommitted changes</td>
    <td colspan="2" align="center">git status</td>
  </tr>
  <tr>
    <td rowspan="3">2.</td>
    <td rowspan="3">pull the feature branch</td>
    <td rowspan="2" colspan="2" align="center">git fetch</td>
  </tr>
  <tr></tr>
  <tr>
    <td>git rebase origin/feature</td>
    <td>git merge origin/feature</td>
  </tr>
  <tr>
    <td>3.</td>
    <td>check out the main branch</td>
    <td colspan="2" align="center">git checkout master</td>
  </tr>
  <tr>
    <td>4.</td>
    <td>pull the main branch</td>
    <td>git rebase origin/master</td>
    <td>git merge origin/master</td>
  </tr>
  <tr>
    <td>5.</td>
    <td>merge feature into main</td>
    <td colspan="2" align="center">git merge --squash feature
  </tr>
  <tr>
    <td>6.</td>
    <td>push the updated master</td>
    <td colspan="2" align="center">git push</td>
  </tr>
  <tr>
    <td>7.</td>
    <td>delete feature locally</td>
    <td colspan="2" align="center">git branch -d feature</td>
  </tr>
  <tr>
    <td>8.</td>
    <td>delete feature from the repo</td>
    <td colspan="2" align="center">git push origin :feature
  </tr>
</table>


### git kill
Safely deletes a git branch.

* never deletes the main branch
* if the branch has unmerged commits, asks the user for confirmation\*
* deletes the given branch from the local machine as well as the repo


## Installation

Git Town is 100% bash script, so it runs anywhere where Git and Bash runs:
OS X, Linux, BSD, and even Windows with Cygwin or something similar.
Installation on OS X can be done using <a href="http://brew.sh" target="_blank">Homebrew</a>,
other platforms need to install manually.

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
      cd [directory of your Git Town clone]
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



## Configuration\*

Git Town operates under the following assumptions:

* You have a **main branch** (typically "development" or "master") from which feature branches are cut, and into which they are merged. In this documentation we will use "master".
* You follow a per-project strategy that prefers either _rebases_ or _merges_ for updating branches, and _squash merges_ or _normal merges_ for merging feature branches into the main branch\*.
* You use a central code repository like [Github](http://github.com) (called __repo__ from now on).

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


## Develop your own scripts

* run all tests: `cucumber`
* run a single test: `cucumber -n 'scenario or feature name`


Some background on the code structure:
* Due to limitations of Bash Script, the functions take normal arguments, and return their result as global variables.
* Each function does the thing it says in a robust way. The "pull_feature_branch" function for example switches to the current feature branch, and then pulls it.



## Release Notes

### 0.3
* <a href="http://cukes.info" target="_blank">Cucumber</a> feature specs
* completely uses local Git repos for testing: https://github.com/Originate/git-town/issues/25
* new configuration file name: .gittownrc instead of the old .main_branch_name
* always cleans up abort and continue scripts
* only makes one fetch from the central repo per session
* specs no longer commit the Git Town configuration file to the repo
* automatically prunes remote branches when fetching updates


### 0.2.2
* fixes "unary" error messages
* lots of output and documentation improvements


### 0.2.1
* better terminal output
* Travis CI improvements
* better documentation

### 0.2
* displays the duration of specs
* only pulls the main branch if it has a remote
* --abort options to abort failed Git Town operations
* --continue options to continue some Git Town operations after fixing the underlying issues
* can be installed through Homebrew
* colored test output
* display summary after tests
* exit with proper status codes
* better documentation


### 0.1
* git hack, git sync, git extract, git ship, git kill
* basic test framework
* Travis CI integration
* self-hosting: uses Git Town for Git Town development
