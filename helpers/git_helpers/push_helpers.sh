#!/bin/bash


# Returns whether the current branch has local updates
# that haven't been pushed to the remote yet.
# Assumes the current branch has a tracking branch
function needs_push {
  if [ "$(git status | grep -c "Your branch is ahead of")" != 0 ]; then
    echo true
  else
    echo false
  fi
}


# Pushes the branch with the given name to origin
function push {
  local current_branch_name=$(get_current_branch_name)
  if [ "$(has_tracking_branch "$current_branch_name")" == true ]; then
    if [ "$(needs_push)" == true ]; then
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
