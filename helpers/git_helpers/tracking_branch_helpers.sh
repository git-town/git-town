#!/bin/bash


# Determines whether the given branch has a remote tracking branch.
function has_tracking_branch {
  local branch_name=$1
  if [ "$(git branch -r | tr -d ' ' | grep -c "^origin\/$branch_name\$")" == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Merges the tracking branch, if one exists, into the current branch
function merge_tracking_branch {
  local branch_name=$(get_current_branch_name)
  if [ "$(has_tracking_branch "$branch_name")" == true ]; then
    merge "origin/$branch_name"
  fi
}


# Merges the tracking branch, if one exists, into the current branch
function rebase_tracking_branch {
  local branch_name=$(get_current_branch_name)
  if [ "$(has_tracking_branch "$branch_name")" == true ]; then
    rebase "origin/$branch_name"
  fi
}
