#!/bin/bash

# Helper methods for dealing with abort/continue scripts


function add_to_abort_script {
  add_to_file "$1" "$abort_script_filename"
}


function add_to_command_list {
  add_to_file "$1" "$command_list_filename"
}


function add_to_file {
  local content=$1
  local filename=$2
  local operator=">"
  if [ -e "$filename" ]; then operator=">>"; fi
  eval "echo '$content' $operator $filename"
}


function add_to_undo_script {
  add_to_file "$1" "$undo_script_filename"
}


function execute_command_list {
  while [ "$(wc -l < "$command_list_filename" | tr -d ' ')" -gt 0 ]; do
    local cmd=$(pop_line "$command_list_filename")
    $cmd
    if [ $? != 0 ]; then exit_with_error; fi
  done

  remove_scripts
}


function exit_with_script_messages {
  local cmd="${program/-/ }"

  echo
  if [ "$(has_script "$abort_script_filename")" == true ]; then
    echo_red "To abort, run \"$cmd --abort\"."
  fi
  if [ "$(has_script "$command_list_filename")" == true ]; then
    echo_red "To continue after you have resolved the conflicts, run \"$cmd --continue\"."
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


function pop_line {
  local file=$1
  local temp=$(temp_filename)
  head -n 1 "$file"
  tail -n +2 "$file" > "$temp"
  mv "$temp" "$file"
}


function prepend_to_command_list {
  if [ "$(has_script "$command_list_filename")" == true ]; then
    local file=$(temp_filename)
    echo "$1" | cat - "$command_list_filename" > "$file" && mv "$file" "$command_list_filename"
  fi
}


function remove_scripts {
  if [ "$(has_script "$abort_script_filename")" == true ]; then
    rm "$abort_script_filename"
  fi
  if [ "$(has_script "$command_list_filename")" == true ]; then
    rm "$command_list_filename"
  fi
  if [ "$(has_script "$undo_script_filename")" == true ]; then
    rm "$undo_script_filename"
  fi
}


function run_abort_script {
  if [ "$(has_script "$abort_script_filename")" == true ]; then
    source "$abort_script_filename"
    remove_scripts
  else
    echo_red "Cannot find abort definition file"
  fi
}


function run_undo_script {
  if [ "$(has_script "$undo_script_filename")" == true ]; then
    source "$undo_script_filename"
    remove_scripts
  else
    echo_red "Cannot find undo definition file"
  fi
}


function write_conflict_abort_script {
  add_to_abort_script "checkout_branch $initial_branch_name"
  if [ "$initial_open_changes" = true ]; then
    add_to_abort_script "restore_open_changes"
  fi
}


function write_merge_conflict_abort_script {
  add_to_abort_script "abort_merge"
  write_conflict_abort_script
}


function write_rebase_conflict_abort_script {
  add_to_abort_script "abort_rebase"
  write_conflict_abort_script
}
