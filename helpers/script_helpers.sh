#!/bin/bash

# Helper methods for dealing with abort/continue scripts

function abort_command {
  local cmd=$(peek_line "$steps")
  eval "abort_$cmd"
  undo_command
}


function continue_command {
  local cmd=$(pop_line "$steps")
  eval "continue_$cmd"
  run_steps "$steps" cleanup
}


function undo_command {
  run_steps "$undo_steps" cleanup
}


function ensure_abortable {
  if [ "$(has_file "$steps")" = false ]; then
    echo_red "Cannot abort"
    exit_with_error
  fi
}


function ensure_continuable {
  if [ "$(has_file "$steps")" = false ]; then
    echo_red "Cannot continue"
    exit_with_error
  fi
}


# Placeholder for any scripts that have no preconditions
function preconditions {
  true
}


function ensure_undoable {
  if [ "$(has_file "$undo_steps")" = false ]; then
    echo_red "Cannot undo"
    exit_with_error
  fi
}


function exit_with_messages {
  if [ "$(has_file "$steps")" = true ]; then
    echo
    echo_red "To abort, run \"$git_command --abort\"."
    echo_red "To continue after you have resolved the conflicts, run \"$git_command --continue\"."
    exit_with_error
  fi
}


function has_lines {
  local file="$1"

  if [ "$(has_file "$file")" = true ] && [ "$(number_of_lines "$file")" -gt 0 ]; then
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


function number_of_lines {
  local file="$1"
  wc -l < "$file" | tr -d ' '
}


function peek_line {
  local file="$1"
  head -n 1 "$file"
}


function pop_line {
  local file="$1"
  peek_line "$file"
  remove_line "$file"
}


function remove_line {
  local file="$1"
  if [ "$(number_of_lines "$file")" -gt 1 ]; then
    local temp=$(temp_filename)
    tail -n +2 "$file" > "$temp"
    mv "$temp" "$file"
  else
    rm "$file"
  fi
}


function remove_step_files {
  if [ "$(has_file "$steps")" = true ]; then
    rm "$steps"
  fi
  if [ "$(has_file "$undo_steps")" = true ]; then
    rm "$undo_steps"
  fi
}


function run {
  if [ "$1" = "--abort" ]; then
    ensure_abortable
    abort_command
  elif [ "$1" = "--undo" ]; then
    ensure_undoable
    undo_command
  elif [ "$1" = "--continue" ]; then
    ensure_continuable
    ensure_no_conflicts
    continue_command
  else
    remove_step_files
    preconditions "$@"
    add_steps
    run_steps "$steps" undoable
  fi

  exit_with_success
}


function run_steps {
  local file="$1"
  local option="$2"

  while [ "$(has_lines "$file")" = true ]; do
    local cmd=$(peek_line "$file")
    if [ "$option" = undoable ]; then
      local undo_cmds=$(undo_commands_for "$cmd")
    fi
    eval "$cmd"

    if [ $? != 0 ]; then
      exit_with_messages
    else
      if [ "$option" = undoable ]; then
        for undo_cmd in "${undo_cmds[@]}"; do
          add_undo_step "$undo_cmd"
        done
      fi
      remove_line "$file"
    fi
  done

  if [ "$option" = cleanup ]; then
    remove_step_files
  fi
}
