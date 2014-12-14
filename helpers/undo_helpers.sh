#!/bin/bash

function undo_steps_for {
  local step_with_arguments
  read -a step_with_arguments <<< "$1" # Split string into array

  local step="${step_with_arguments[0]}"
  local arguments="${step_with_arguments[*]:1}"

  local fn="undo_steps_for_$step"

  if [ "$(type "$fn" 2>&1 | grep -c 'not found')" = 0 ]; then
    eval "$fn $arguments"
  fi
}




function undo_steps_for_checkout {
  local branch=$(get_current_branch_name)
  echo "checkout $branch"
}


function undo_steps_for_checkout_main_branch {
  undo_steps_for_checkout
}


function undo_steps_for_commit_open_changes {
  local branch=$(get_current_branch_name)
  local sha=$(sha_of_branch "$branch")
  echo "reset_to_sha $sha"
  if [ "$(has_tracking_branch "$branch")" = true ]; then
    echo "push_branch $branch force"
  fi
}


function undo_steps_for_create_and_checkout_feature_branch {
  local branch=$(get_current_branch_name)
  local branch_to_create="$1"
  echo "checkout $branch"
  echo "delete_branch $branch_to_create"
}


function undo_steps_for_delete_branch {
  local branch_to_delete="$1"
  local sha=$(sha_of_branch "$branch_to_delete")
  echo "create_branch $branch_to_delete $sha"
  if [ "$(has_tracking_branch "$branch_to_delete")" = true ]; then
    echo "push_branch $branch_to_delete"
  fi
}


function undo_steps_for_stash_open_changes {
  echo "restore_open_changes"
}
