#!/bin/bash


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


# Push and pull the current branch
function sync_branch {
  local strategy=$1
  pull_branch "$strategy"
  push_branch
}
