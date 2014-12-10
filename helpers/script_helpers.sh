#!/bin/bash

# Helper methods for dealing with abort/continue scripts

function abort_command {
  local cmd=$(peek_command "$command_list")
  eval "abort_$cmd"
  undo_command
}


function continue_command {
  local cmd=$(peek_command "$command_list")
  pop_command "$command_list"
  eval "continue_$cmd"
  run_command_list "$command_list" 'cleanup'
}


function undo_command {
  run_command_list "$undo_command_list" 'cleanup'
}


function ensure_abortable {
  if [ "$(has_file "$command_list")" = false ]; then
    echo_red "Cannot abort"
    exit_with_error
  fi
}


function ensure_continuable {
  if [ "$(has_file "$command_list")" = false ]; then
    echo_red "Cannot continue"
    exit_with_error
  fi
}


function ensure_undoable {
  if [ "$(has_file "$undo_command_list")" = false ]; then
    echo_red "Cannot undo"
    exit_with_error
  fi
}


function exit_with_messages {
  if [ "$(has_file "$command_list")" = true ]; then
    echo
    echo_red "To abort, run \"$git_command --abort\"."
    echo_red "To continue after you have resolved the conflicts, run \"$git_command --continue\"."
    exit_with_error
  fi
}


function has_commands {
  local file="$1"

  if [ "$(has_file "$file")" = true ] && [ "$(number_of_commands "$file")" -gt 0 ]; then
    echo true
  else
    echo false
  fi
}


function has_file {
  local file="$1"

  if [ -n "$file" -a -f "$file" ]; then
    echo true
  else
    echo false
  fi
}


function number_of_commands {
  local file="$1"
  wc -l < "$file" | tr -d ' '
}


function peek_command {
  local file="$1"
  head -n 1 "$file"
}


function pop_command {
  local file="$1"
  if [ "$(number_of_commands "$file")" -gt 1 ]; then
    local temp=$(temp_filename)
    tail -n +2 "$file" > "$temp"
    mv "$temp" "$file"
  else
    rm "$file"
  fi
}


function remove_command_lists {
  if [ "$(has_file "$abort_command_list")" = true ]; then
    rm "$abort_command_list"
  fi
  if [ "$(has_file "$command_list")" = true ]; then
    rm "$command_list"
  fi
  if [ "$(has_file "$undo_command_list")" = true ]; then
    rm "$undo_command_list"
  fi
}


function run {
  if [ "$1" = "--abort" ]; then
    ensure_abortable
    abort_command
    remove_command_lists
  elif [ "$1" = "--undo" ]; then
    ensure_undoable
    undo_command
  elif [ "$1" = "--continue" ]; then
    ensure_continuable
    ensure_no_conflicts
    continue_command
    remove_command_lists
  else
    remove_command_lists
    build_command_list "$@"
    run_command_list "$command_list" 'build_undo'
  fi

  exit_with_success
}


function run_command_list {
  local file="$1"
  local option="$2"

  while [ "$(has_commands "$file")" = true ]; do
    local cmd=$(peek_command "$file")
    if [ "$option" = 'build_undo' ]; then
      local undo_cmds=$(undo_commands_for "$cmd")
    fi
    eval "$cmd"

    if [ $? != 0 ]; then
      exit_with_messages
    else
      if [ "$option" = 'build_undo' ]; then
        for undo_cmd in "${undo_cmds[@]}"; do
          add_to_undo_command_list "$undo_cmd"
        done
      fi
      pop_command "$file"
    fi
  done

  if [ "$option" = 'cleanup' ]; then
    remove_command_lists
  fi
}
