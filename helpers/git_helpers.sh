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


# Returns whether or not the two branches are in sync
function branches_in_sync  {
  local branch_one=$1
  local branch_two=$2
  if [ "$(git rev-list --left-right "$branch_one..$branch_two" | wc -l | tr -d ' ')" = 0 ]; then
    echo true
  else
    echo false
  fi
}


# Returns whether or not the branch is in sync with its tracking branch
function branch_in_sync {
  local branch_name=$1
  if [ "$(has_tracking_branch "$branch_name")" = true ]; then
    local tracking_branch_name=$(tracking_branch "$branch_name")
    if [ "$(branches_in_sync "$branch_name" "$tracking_branch_name")" = true ]; then
      echo true
    else
      echo false
    fi
  else
    echo true
  fi
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


# Exits the application with an error message if the
# branch is not in sync
function ensure_branch_in_sync {
  local branch_name=$1
  if [ "$(branch_in_sync "$branch_name")" == false ]; then
    echo_error_header
    echo_error "The '$branch_name' branch is out of sync. Run 'git sync' to resolve."
    exit_with_error
  fi
}


# Exits the application with an error message if the
# repository has a branch with the given name.
function ensure_does_not_have_branch {
  local branch_name=$1
  if [ "$(has_branch "$branch_name")" = true ]; then
    echo_error_header
    echo_error "A branch named '$branch_name' already exists"
    exit_with_error
  fi
}



# Exits the application with an error message if the
# repository does not have a branch with the given name.
function ensure_has_branch {
  local branch_name=$1
  if [ "$(has_branch "$branch_name")" == false ]; then
    echo_error_header
    echo_error "There is no branch named '$branch_name'."
    exit_with_error
  fi
}


# Exit if the given branch does not have shippable changes
function ensure_has_shippable_changes {
  local branch_name=$1
  if [ "$(has_shippable_changes "$branch_name")" == false ]; then
    return_to_initial_branch

    echo_error_header
    echo_error "The branch '$branch_name' has no shippable changes."
    exit_with_error
  fi
}


# Exits the application with an error message if the supplied branch is
# not a feature branch
function ensure_is_feature_branch {
  local branch_name=$1
  local error_message=$2
  if [ "$(is_feature_branch "$branch_name")" == false ]; then
    error_is_not_feature_branch

    echo_error_header
    echo_error "The branch '$branch_name' is not a feature branch. $error_message"
    exit_with_error
  fi
}


# Exits with an error message if there are unresolved conflicts
function ensure_no_conflicts {
  if [ "$(has_conflicts)" == true ]; then
    echo_error_header
    echo_error "$*"
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


# Exits the application with an error message if the current branch is
# not a feature branch
function ensure_on_feature_branch {
  local error_message="$*"
  local branch_name=$(get_current_branch_name)
  ensure_is_feature_branch "$branch_name" "$error_message"
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
  if [ "$(git branch | tr -d '* ' | grep -c "^$branch_name\$")" = 0 ]; then
    echo false
  else
    echo true
  fi
}


# Returns true if there are conflicts
function has_conflicts {
  if [ "$(git status | grep -c 'Unmerged paths')" == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether there are open changes in Git.
function has_open_changes {
  if [ "$(git status --porcelain | wc -l | tr -d ' ')" == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether the given branch has shippable changes
function has_shippable_changes {
  local branch_name=$1
  if [ "$(git diff --quiet "$main_branch_name..$branch_name" ; echo $?)" == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether the given branch has a remote tracking branch.
function has_tracking_branch {
  local branch_name=$1
  if [ -z "$(tracking_branch "$branch_name")" ]; then
    echo false
  else
    echo true
  fi
}


# Returns true if the current branch is a feature branch
function is_feature_branch {
  local branch_name=$1
  if [ "$branch_name" == "$main_branch_name" -o "$(echo "$non_feature_branch_names" | tr ',' '\n' | grep -c "$branch_name")" == 1 ]; then
    echo false
  else
    echo true
  fi
}


# Returns the names of local branches that have been merged into main
function local_merged_branches {
  git branch --merged "$main_branch_name" | tr -d ' ' | sed 's/\*//g'
}


# Merges the given branch into the current branch
function merge_branch {
  local branch_name=$1
  local current_branch_name=$(get_current_branch_name)
  run_command "git merge --no-edit $branch_name"
  if [ $? != 0 ]; then error_merge_branch; fi
}


# Returns whether the current branch has local updates
# that haven't been pushed to the remote yet.
# Assumes the current branch has a tracking branch
function needs_pushing {
  if [ "$(git status | grep -c "Your branch is ahead of")" != 0 ]; then
    echo true
  else
    echo false
  fi
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


# Determines whether the current branch has a rebase in progress
function rebase_in_progress {
  if [ "$(git status | grep -c "rebase in progress")" == 1 ]; then
    echo true
  else
    echo false
  fi
}


# Returns the names of remote branches that have been merged into main
function remote_merged_branches {
  git branch -r --merged "$main_branch_name" | grep -v HEAD | tr -d ' ' | sed 's/origin\///g'
}


# Returns the url for the remote
function remote_url {
  git remote -v | grep "origin.*fetch" | awk '{print $2}'
}


# Returns the domain of the remote repository
function remote_domain {
  remote_url | sed -E "s#(https?://([^@]*@)?|git@)([^/:]+).*#\3#"
}


# Returns the USER/REPO for the remote repository
function remote_repository_name {
  remote_url | sed "s#.*[:/]\([^/]*/[^/]*\)\.git#\1#"
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


# Returns the SHA that the given branch points to
function sha_of_branch {
  local branch_name=$1
  git rev-parse "$branch_name"
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


# Returns the name of the tracking branch for the given branch
function tracking_branch {
  local branch_name=$1
  local result
  result=$(git rev-parse --abbrev-ref "${branch_name}@{u}")
  if [ "$?" = 0 ]; then echo "$result"; fi
}
