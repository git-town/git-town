#!/usr/bin/env bash

# Helper methods for dealing with abort/continue scripts


function abort_command {
  local cmd=$(peek_line "$steps_file")
  eval "abort_$cmd"
  undo_command
}


function continue_command {
  local cmd=$(pop_line "$steps_file")
  eval "continue_$cmd"
  run_steps "$steps_file" undoable
}


function skip_command {
  local cmd=$(pop_line "$steps_file")
  eval "abort_$cmd"
  undo_current_branch_steps
  skip_current_branch_steps
  run_steps "$steps_file" undoable
}


function undo_command {
  run_steps "$undo_steps_file" cleanup
}


function ensure_abortable {
  if [ "$(has_file "$steps_file")" = false ]; then
    echo_red "Cannot abort"
    exit_with_error
  fi
}


function ensure_continuable {
  if [ "$(has_file "$steps_file")" = false ]; then
    echo_red "Cannot continue"
    exit_with_error
  fi
}


function ensure_skippable {
  if [ "$(has_file "$steps_file")" = false ]; then
    echo_red "Cannot skip"
    exit_with_error
  fi
}


function ensure_undoable {
  if [ "$(has_file "$undo_steps_file")" = false ]; then
    echo_red "Cannot undo"
    exit_with_error
  fi
}


function exit_with_messages {
  if [ "$(has_file "$steps_file")" = true ]; then
    echo
    echo_red "To abort, run \"$git_command --abort\"."
    echo_red "To continue after you have resolved the conflicts, run \"$git_command --continue\"."
    if [ "$(skippable)" = true ]; then
      echo_red "$(skip_message), run \"$git_command --skip\"."
    fi
    exit_with_error
  fi
}


# Placeholder for any scripts that have no preconditions
function preconditions {
  true
}


function remove_step_files {
  if [ "$(has_file "$steps_file")" = true ]; then
    rm "$steps_file"
  fi
  if [ "$(has_file "$undo_steps_file")" = true ]; then
    rm "$undo_steps_file"
  fi
}


function run {
  if [ "$1" = "--abort" ]; then
    ensure_abortable
    abort_command
  elif [ "$1" = "--continue" ]; then
    ensure_continuable
    ensure_no_conflicts
    continue_command
  elif [ "$1" = "--skip" ]; then
    ensure_skippable
    skip_command
  elif [ "$1" = "--undo" ]; then
    ensure_undoable
    undo_command
  else
    remove_step_files
    preconditions "$@"
    steps > "$steps_file"
    run_steps "$steps_file" undoable
  fi

  exit_with_success
}


# possible values for option
#   undoable - builds an undo_steps_file
#   cleanup - calls remove_step_files after successfully running all steps
function run_steps {
  local file="$1"
  local option="$2"

  while [ "$(has_lines "$file")" = true ]; do
    local step=$(peek_line "$file")
    if [ "$option" = undoable ]; then
      local undo_steps=$(undo_steps_for "$step")
    fi
    eval "$step"

    if [ $? != 0 ]; then
      exit_with_messages
    else
      if [ "$option" = undoable ]; then
        add_undo_steps "$undo_steps"
      fi
      remove_line "$file"
    fi
  done

  if [ "$option" = cleanup ]; then
    remove_step_files
  fi
}


function skip_current_branch_steps {
  while [ "$(has_lines "$steps_file")" = true ]; do
    if [[ "$(peek_line "$steps_file")" =~ ^checkout ]]; then
      break
    else
      remove_line "$steps_file"
    fi
  done
}


function undo_current_branch_steps {
  while [ "$(has_lines "$undo_steps_file")" = true ]; do
    local step=$(peek_line "$undo_steps_file")
    if [[ "$step" =~ ^checkout ]]; then
      break
    else
      eval "$step"
      remove_line "$undo_steps_file"
    fi
  done
}


# Placeholder for any scripts that do have the skip interface
function skippable {
  echo false
}
