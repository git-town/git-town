#!/usr/bin/env bash


function add_undo_steps {
  local steps=$1

  prepend_to_file "$steps" "$UNDO_STEPS_FILE"
}


# Returns the post_undo steps for the given command
function post_undo_steps_for {
  local command=$1
  _undo_steps "$command" "post_undo_steps_for_"
}


# Returns the undo steps for the given command
function undo_steps_for {
  local command=$1
  _undo_steps "$command" "undo_steps_for_"
}


function _undo_steps {
  local command=$1
  local step_name=$2

  local step_with_arguments
  read -a step_with_arguments <<< "$command" # Split string into array
  local step="${step_with_arguments[0]}"
  local arguments="${step_with_arguments[*]:1}"
  local fn="$step_name$step"
  if type "${fn}" > /dev/null 2>&1; then
    eval "$fn $arguments"
  fi
}
