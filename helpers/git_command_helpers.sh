#!/bin/bash

# Helper methods for running git commands


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
  local SHAs="$*"
  run_command "git cherry-pick $SHAs"
  if [ $? != 0 ]; then error_cherry_pick "$SHAs"; fi
}


# Commits all open changes into the current branch
function commit_open_changes {
  run_command "git add -A"
  run_command "git commit -m 'WIP on $(get_current_branch_name)'"
}


# Continues merge if one is in progress
function continue_merge {
  if [ "$(has_open_changes)" == true ]; then
    run_command "git commit --no-edit"
  fi
}


# Continues rebase if one is in progress
function continue_rebase {
  if [ "$(rebase_in_progress)" == true ]; then
    run_command "git rebase --continue"
  fi
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


# Discard open changes
function discard_open_changes {
  run_command 'git reset --hard'
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


# Merges the given branch into the current branch
function merge_branch {
  local branch_name=$1
  local current_branch_name=$(get_current_branch_name)
  run_command "git merge --no-edit $branch_name"
  if [ $? != 0 ]; then error_merge_branch; fi
}


# Pulls updates of the feature branch from the remote repo
function pull_branch {
  local strategy=$1
  local current_branch_name=$(get_current_branch_name)
  if [ -z "$strategy" ]; then strategy='merge'; fi
  if [ "$(has_tracking_branch "$current_branch_name")" == true ]; then
    fetch_repo
    run_command "git $strategy origin/$current_branch_name"
    if [ $? != 0 ]; then error_pull_branch "$current_branch_name"; fi
  fi
}


# Pulls updates of the current branch from the upstream repo
function pull_upstream_branch {
  local current_branch_name=$(get_current_branch_name)
  fetch_upstream
  run_command "git rebase upstream/$current_branch_name"
  if [ $? != 0 ]; then error_pull_upstream_branch; fi
}


# Pushes the branch with the given name to origin
function push_branch {
  local current_branch_name=$(get_current_branch_name)
  if [ "$(has_tracking_branch "$current_branch_name")" == true ]; then
    if [ "$(needs_pushing)" == true ]; then
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


# Resets the current branch to the commit described by the given SHA
function reset_to_sha {
  local sha=$1
  run_command "git reset $sha"
}


# Unstashes changes that were stashed in the beginning of a script.
#
# Only does this if there were open changes when the script was started.
function restore_open_changes {
  if [ "$initial_open_changes" = true ]; then
    run_command "git stash pop"
  fi
}


# Squash merges the given branch into the current branch
function squash_merge {
  local branch_name=$1
  local commit_message=$2
  local current_branch_name=$(get_current_branch_name)
  run_command "git merge --squash $branch_name"
  if [ $? != 0 ]; then error_squash_merge; fi
  if [ "$commit_message" == "" ]; then
    run_command "git commit -a"
  else
    run_command "git commit -a -m '$commit_message'"
  fi
  if [ $? != 0 ]; then error_empty_commit; fi
}


# Stashes uncommitted changes if they exist.
function stash_open_changes {
  if [ "$initial_open_changes" = true ]; then
    run_command "git stash -u"
  fi
}


# Push and pull the current branch
function sync_branch {
  local strategy=$1
  pull_branch "$strategy"
  push_branch
}
