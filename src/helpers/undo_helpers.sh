#!/usr/bin/env bash


function add_undo_steps {
  local steps=$1

  prepend_to_file "$steps" "$UNDO_STEPS_FILE"
}


# Returns the post_undo steps for the given command
function post_undo_steps_for {
  local command=$1

  local step_with_arguments
  read -a step_with_arguments <<< "$command" # Split string into array

  local step="${step_with_arguments[0]}"
  local arguments="${step_with_arguments[*]:1}"
  local fn="post_undo_steps_for_$step"
  if [ "$(type "$fn" 2>&1 | grep -c 'not found')" = 0 ]; then
    eval "$fn $arguments"
  fi
}


# Returns the undo steps for the given command
function undo_steps_for {
  local step_with_arguments
  read -a step_with_arguments <<< "$1" # Split string into array

  local step="${step_with_arguments[0]}"
  local arguments="${step_with_arguments[*]:1}"

  local fn="undo_steps_for_$step"

  if [ "$(type "$fn" 2>&1 | grep -c 'not found')" = 0 ]; then
    eval "$fn $arguments"
  fi
}
