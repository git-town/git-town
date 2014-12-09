#!/bin/bash

function add_undo_command_for {
  local full_cmd
  read -a full_cmd <<< "$1"

  local cmd="${full_cmd[0]}"

  if [ "$cmd" = "stash_open_changes" ]; then
    add_to_undo_command_list "restore_open_changes"
  fi
}
