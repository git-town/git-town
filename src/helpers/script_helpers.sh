#!/usr/bin/env bash

# Helper methods for dealing with abort/continue scripts


function abort_command {
  local cmd=$(peek_line "$STEPS_FILE")
  eval "abort_$cmd"
  undo_command
}


function continue_command {
  local cmd=$(pop_line "$STEPS_FILE")
  eval "continue_$cmd"
  run_steps "$STEPS_FILE" undoable
}


function skip_command {
  local cmd=$(pop_line "$STEPS_FILE")
  eval "abort_$cmd"
  undo_current_branch_steps
  skip_current_branch_steps "$STEPS_FILE"
  run_steps "$STEPS_FILE" undoable
}


function undo_command {
  run_steps "$UNDO_STEPS_FILE" cleanup
}


function ensure_abortable {
  if [ "$(has_file "$STEPS_FILE")" = false ]; then
    echo_red "Cannot abort"
    exit_with_error
  fi
}


function ensure_continuable {
  if [ "$(has_file "$STEPS_FILE")" = false ]; then
    echo_red "Cannot continue"
    exit_with_error
  fi

  if [ "$(has_lines "$STEPS_FILE")" = false ]; then
    echo_red "The last command finished successfully and cannot be continued"
    exit_with_error
  fi
}


function ensure_skippable {
  if [ "$(has_file "$STEPS_FILE")" = false ]; then
    echo_red "Cannot skip"
    exit_with_error
  fi
}


function ensure_undoable {
  if [ "$(has_file "$UNDO_STEPS_FILE")" = false ]; then
    echo_red "Nothing to undo"
    exit_with_error
  fi
}


function exit_with_messages {
  if [ "$(has_file "$STEPS_FILE")" = true ]; then
    echo
    echo_red "To abort, run \"$GIT_COMMAND --abort\"."
    echo_red "To continue after you have resolved the conflicts, run \"$GIT_COMMAND --continue\"."
    if [ "$(skippable)" = true ]; then
      echo_red "$(skip_message_prefix), run \"$GIT_COMMAND --skip\"."
    fi
    exit_with_error newline
  fi
}


function get_expected_previous_branch {
  if [ "$(has_branch "$INITIAL_PREVIOUS_BRANCH_NAME")" = true ]; then

    # current branch is unchanged
    if [ "$(get_current_branch_name)" = "$INITIAL_BRANCH_NAME" ]; then
      echo "$INITIAL_PREVIOUS_BRANCH_NAME"

    # current branch is deleted
    elif [ "$(has_branch "$INITIAL_BRANCH_NAME")" = false ]; then
      echo "$INITIAL_PREVIOUS_BRANCH_NAME"

    # current branch is new
    else
      echo "$INITIAL_BRANCH_NAME"
    fi

  # previous branch is deleted
  else
    echo "$MAIN_BRANCH_NAME"
  fi
}


# Parses arguments, validates necessary state
# This should be overriden by commands when necessary
function preconditions {
  true
}


function remove_step_files {
  if [ "$(has_file "$STEPS_FILE")" = true ]; then
    rm "$STEPS_FILE"
  fi
  if [ "$(has_file "$UNDO_STEPS_FILE")" = true ]; then
    rm "$UNDO_STEPS_FILE"
  fi
}


function restore_proper_previous_branch {
  set_previous_branch "$(get_expected_previous_branch)"
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
    steps > "$STEPS_FILE"
    run_steps "$STEPS_FILE" undoable

    # Apply global previous-branch restoration
    # only if we're at the git root.
    # Any branch restoration logic for the other
    # case needs to be handled in a custom way by
    # the command itself
    if [ "$IN_SUB_FOLDER" = false ]; then
      restore_proper_previous_branch
    fi
  fi
}


# possible values for option
#   undoable - builds an UNDO_STEPS_FILE
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

  echo # trailing newline (each git command prints a leading newline)
}


# Skip any steps on the current branch
function skip_current_branch_steps {
  local file="$1"

  while [ "$(has_lines "$file")" = true ]; do
    if [[ "$(peek_line "$file")" =~ ^checkout ]]; then
      # Add empty line to file to ensure no step if lost
      prepend_to_file "" "$file"
      break
    else
      remove_line "$file"
    fi
  done
}


# Returns whether or not the current step can be skipped
# This should be overriden in commands when necessary
function skippable {
  echo false
}


# Undo any steps on the current branch
function undo_current_branch_steps {
  while [ "$(has_lines "$UNDO_STEPS_FILE")" = true ]; do
    local step=$(peek_line "$UNDO_STEPS_FILE")
    if [[ "$step" =~ ^checkout ]]; then
      break
    else
      eval "$step"
      remove_line "$UNDO_STEPS_FILE"
    fi
  done
}
