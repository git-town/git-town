#!/bin/bash

function add_undo_command_for {
  local branch=$1
  local full_cmd
  read -a full_cmd <<< "$2"

  local cmd="${full_cmd[0]}"

  if [ "$cmd" = "stash_open_changes" ]; then
    add_to_undo_command_list "restore_open_changes"
  elif [ "$cmd" = "checkout" ] || [ "$cmd" = "checkout_main_branch" ]; then
    add_to_undo_command_list "checkout $branch"
  elif [ "$cmd" = "create_and_checkout_feature_branch" ]; then
    local feature_branch="${full_cmd[1]}"
    add_to_undo_command_list "delete_branch $feature_branch"
    add_to_undo_command_list "checkout $branch"
  fi
}
