#!/bin/bash

# Helper methods for working with Git.


# Abort a cherry-pick
function abort_cherry_pick {
  run_command "git cherry-pick --abort"
}


# Abort a merge
function abort_merge {
  run_command "git merge --abort"
}


# Abort a rebase
function abort_rebase {
  run_command "git rebase --abort"
}


# Checks out the branch with the given name.
#
# Skips this operation if the requested branch
# is already checked out.
function checkout_branch {
  local branch_name=$1
  if [ ! "$(get_current_branch_name)" = "$branch_name" ]; then
    run_command "git checkout $branch_name"
  fi
}


# Checks out the main development branch in Git.
#
# Skips the operation if we already are on that branch.
function checkout_main_branch {
  checkout_branch "$main_branch_name"
}


# Cherry picks the SHAs into the current branch
function cherry_pick {
  local SHAs=$*
  run_command "git cherry-pick $SHAs"
  if [ $? != 0 ]; then error_cherry_pick "$SHAs"; fi
}


# Commits all open changes into the current branch
function commit_open_changes {
  run_command "git add -A"
  run_command "git commit -m 'WIP on $(get_current_branch_name)'"
}


# Cuts a new branch off the given parent branch, and checks it out.
function create_and_checkout_branch {
  local new_branch_name=$1
  local parent_branch_name=$2
  run_command "git checkout -b $new_branch_name $parent_branch_name"
}


# Creates a new feature branch with the given name.
#
# The feature branch is cut off the main development branch.
function create_and_checkout_feature_branch {
  create_and_checkout_branch "$1" "$main_branch_name"
}


# Deletes the given branch from both the local machine and on remote.
function delete_branch {
  local branch_name=$1
  local force=$2
  if [ "$(has_tracking_branch "$branch_name")" == true ]; then
    delete_remote_branch "$branch_name"
  fi
  delete_local_branch "$branch_name" "$force"
}


# Deletes the local branch with the given name
function delete_local_branch {
  local branch_name=$1
  local op="d"
  if [ "$2" == "force" ]; then op="D"; fi
  run_command "git branch -$op $branch_name"
}


# Deletes the remote branch with the given name
function delete_remote_branch {
  local branch_name=$1
  run_command "git push origin :${branch_name}"
}


# Exists the application with an error message if the
# repository does not have a branch with the given name.
function ensure_has_branch {
  local branch_name=$1
  if [ "$(has_branch "$branch_name")" == false ]; then
    echo_error_header
    echo_error "There is no branch named '$branch_name'."
    exit_with_error
  fi
}


# Exists the application with an error message if the
# current working directory contains uncommitted changes.
function ensure_no_open_changes {
  if [ "$(has_open_changes)" == true ]; then
    error_has_open_changes

    echo_error_header
    echo_error "$*"
    exit_with_error
  fi
}


# Exists the application with an error message if the working directory
# is on the main development branch.
function ensure_on_feature_branch {
  local error_message=$1
  local branch_name=$(get_current_branch_name)
  if [ "$(is_feature_branch "$branch_name")" == false ]; then
    error_not_on_feature_branch

    echo_error_header
    echo_error "The branch '$branch_name' is not a feature branch. $error_message"
    exit_with_error
  fi
}


# Called by pull_branch when the merge/rebase fails with conflicts
function error_pull_branch {
  if [ "$(is_feature_branch "$1")" == true ]; then
    error_pull_feature_branch
  elif [ "$1" == "$main_branch_name" ]; then
    error_pull_main_branch
  else
    error_pull_non_feature_branch
  fi
}


# Fetches updates from the central repository.
#
# It is safe to call this method multiple times per session,
# since it makes sure that it fetches updates only once per session
# by tracking this through the global variable $repo_fetched.
function fetch_repo {
  if [ "$repo_fetched" == false ]; then
    run_command "git fetch --prune"
    repo_fetched=true
  fi
}
repo_fetched=false


# Fetches changes from the upstream repository
function fetch_upstream {
  run_command "git fetch upstream"
}


# Returns the current branch name
function get_current_branch_name {
  git rev-parse --abbrev-ref HEAD
}


# Returns true if the repository has a branch with the given name
function has_branch {
  local branch_name=$1
  if [ `git branch | tr -d '* ' | grep "^$branch_name$" | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether there are open changes in Git.
function has_open_changes {
  if [ `git status --porcelain | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether the given branch has a remote tracking branch.
function has_tracking_branch {
  local branch_name=$1
  if [ `git branch -vv | grep "$branch_name" | grep "\[origin\/$branch_name.*\]" | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether the given branch is ahead of main
function is_ahead_of_main {
  local branch_name=$1
  if [ `git log --oneline $main_branch_name..$branch_name | wc -l` == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Returns true if the current branch is a feature branch
function is_feature_branch {
  local branch_name=$1
  if [ "$branch_name" == "$main_branch_name" -o `echo $non_feature_branch_names | tr ',' '\n' | grep $branch_name | wc -l` == 1 ]; then
    echo false
  else
    echo true
  fi
}


# Returns the names of local branches that have been merged into main
function local_merged_branches {
  git branch --merged $main_branch_name | tr -d ' ' | sed 's/\*//g'
}


# Merges the given branch into the current branch
function merge_branch {
  local branch_name=$1
  local current_branch_name=`get_current_branch_name`
  run_command "git merge --no-edit $branch_name"
  if [ $? != 0 ]; then error_merge_branch; fi
}


# Returns whether the current branch has local updates
# that haven't been pushed to the remote yet.
# Assumes the current branch has a tracking branch
function needs_pushing {
  if [ `git status | grep "Your branch is ahead of" | wc -l` != 0 ]; then
    echo true
  else
    echo false
  fi
}


# Pulls updates of the feature branch from the remote repo
function pull_branch {
  local strategy=$1
  local current_branch_name=`get_current_branch_name`
  if [ -z $strategy ]; then strategy='merge'; fi
  if [ `has_tracking_branch $current_branch_name` == true ]; then
    fetch_repo
    run_command "git $strategy origin/$current_branch_name"
    if [ $? != 0 ]; then error_pull_branch $current_branch_name; fi
  fi
}


# Pulls updates of the current branch from the upstream repo
function pull_upstream_branch {
  local current_branch_name=`get_current_branch_name`
  fetch_upstream
  run_command "git rebase upstream/$current_branch_name"
  if [ $? != 0 ]; then error_pull_upstream_branch; fi
}


# Pushes the branch with the given name to origin
function push_branch {
  local current_branch_name=`get_current_branch_name`
  if [ `has_tracking_branch $current_branch_name` == true ]; then
    if [ `needs_pushing` == true ]; then
      run_command "git push"
    fi
  else
    run_command "git push -u origin $current_branch_name"
  fi
}


# Pushes tags to the remote
function push_tags {
  run_command "git push --tags"
}


# Returns the names of remote branches that have been merged into main
function remote_merged_branches {
  git branch -r --merged $main_branch_name | grep -v HEAD | tr -d ' ' | sed 's/origin\///g'
}


# Returns the url for the remote with the specified name
function remote_url {
  git remote -v | grep "$1.*fetch" | awk '{print $2}'
}


# Resets the current branch to the commit described by the given SHA
function reset_to_sha {
  local sha=$1
  run_command 'git reset $sha'
}


# Unstashes changes that were stashed in the beginning of a script.
#
# Only does this if there were open changes when the script was started.
function restore_open_changes {
  if [ $initial_open_changes = true ]; then
    run_command "git stash pop"
  fi
}


# Returns the SHA that the given branch points to
function sha_of_branch {
  local branch_name=$1
  git rev-parse $branch_name
}


# Squash merges the given branch into the current branch
function squash_merge {
  local branch_name=$1
  local commit_message=$2
  local current_branch_name=`get_current_branch_name`
  run_command "git merge --squash $branch_name"
  if [ $? != 0 ]; then error_squash_merge; fi
  if [ "$commit_message" == "" ]; then
    run_command "git commit -a"
  else
    run_command "git commit -a -m '$commit_message'"
  fi
}


# Stashes uncommitted changes if they exist.
function stash_open_changes {
  if [ $initial_open_changes = true ]; then
    run_command "git stash -u"
  fi
}

# Stashes uncommitted changes if they exist.
function sync_main_branch {
  local current_branch_name=`get_current_branch_name`
  checkout_main_branch
  pull_branch 'rebase'
  push_branch
  checkout_branch $current_branch_name
}
