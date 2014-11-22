#!/bin/bash

# Helper methods for dealing with abort/continue scripts


function add_to_abort_script {
  add_to_script $1 "$abort_script_filename"
}


function add_to_continue_script {
  add_to_script "$1" "$continue_script_filename"
}


function add_to_script {
  local content=$1
  local filename=$2
  local operator=">"
  if [ -e "$filename" ]; then operator=">>"; fi
  eval "echo '$content' $operator $filename"
}


function add_to_undo_script {
  add_to_script "$1" "$undo_script_filename"
}


function create_merge_conflict_abort_script {
  add_to_abort_script "initial_open_changes=$initial_open_changes"
  add_to_abort_script "abort_merge"
  add_to_abort_script "checkout_branch $initial_branch_name"
  add_to_abort_script "restore_open_changes"
}


function create_rebase_conflict_abort_script {
  add_to_abort_script "initial_open_changes=$initial_open_changes"
  add_to_abort_script "abort_rebase"
  add_to_abort_script "checkout_branch $initial_branch_name"
  add_to_abort_script "restore_open_changes"
}


function exit_with_abort_continue_messages {
  local cmd=$1

  echo
  if [ "$(has_script "$abort_script_filename")" == true ]; then
    echo_red "To abort, run \"git $cmd --abort\"."
  fi
  if [ "$(has_script "$continue_script_filename")" == true ]; then
    echo_red "To continue after you have resolved the conflicts, run \"git $cmd --continue\"."
  fi
  exit_with_error
}


function has_script {
  if [ -n "$1" -a -f "$1" ]; then
    echo true
  else
    echo false
  fi
}


function remove_abort_continue_scripts {
  if [ "$(has_script "$abort_script_filename")" == true ]; then
    rm "$abort_script_filename"
  fi
  if [ "$(has_script "$continue_script_filename")" == true ]; then
    rm "$continue_script_filename"
  fi
  if [ "$(has_script "$undo_script_filename")" == true ]; then
    rm "$undo_script_filename"
  fi
}


function run_abort_script {
  if [ "$(has_script "$abort_script_filename")" == true ]; then
    source "$abort_script_filename"
    remove_abort_continue_scripts
  else
    echo_red "Cannot find abort definition file"
  fi
}


function run_continue_script {
  if [ "$(has_script "$continue_script_filename")" == true ]; then
    source "$continue_script_filename"
    remove_abort_continue_scripts
  else
    echo_red "Cannot find continue definition file"
  fi
}


function run_undo_script {
  if [ "$(has_script "$undo_script_filename")" == true ]; then
    source "$undo_script_filename"
    remove_abort_continue_scripts
  else
    echo_red "Cannot find undo definition file"
  fi
}
