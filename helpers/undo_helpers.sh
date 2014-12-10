#!/bin/bash

function undo_commands_for {
  local cmd_with_arguments
  read -a cmd_with_arguments <<< "$1"

  local cmd="${cmd_with_arguments[0]}"

  if [ "$cmd" = "stash_open_changes" ]; then
    echo "restore_open_changes"
  elif [ "$cmd" = "checkout" ] || [ "$cmd" = "checkout_main_branch" ]; then
    local branch=$(get_current_branch_name)
    echo "checkout $branch"
  elif [ "$cmd" = "create_and_checkout_feature_branch" ]; then
    local branch=$(get_current_branch_name)
    local branch_to_create="${cmd_with_arguments[1]}"
    echo "checkout $branch"
    echo "delete_branch $branch_to_create"
  elif [ "$cmd" = "delete_branch" ]; then
    local branch_to_delete="${cmd_with_arguments[1]}"
    local sha=$(sha_of_branch "$branch_to_delete")
    echo "create_branch $branch_to_delete $sha"
    if [ "$(has_tracking_branch "$branch_to_delete")" = true ]; then
      echo "push_branch $branch_to_delete"
    fi
  elif [ "$cmd" = "commit_open_changes" ]; then
    local branch=$(get_current_branch_name)
    local sha=$(sha_of_branch "$branch")
    echo "reset_to_sha $sha"
    if [ "$(has_tracking_branch "$branch")" = true ]; then
      echo "push_branch $branch force"
    fi
  fi
}
