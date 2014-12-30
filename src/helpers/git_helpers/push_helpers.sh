#!/usr/bin/env bash


# Returns whether the given branch is in sync with its tracking branch
function needs_push {
  local branch_name=$1
  local tracking_branch_name="origin/$branch_name"
  if [ "$(git rev-list --left-right "$branch_name...$tracking_branch_name" | wc -l | tr -d ' ')" != 0 ]; then
    echo true
  else
    echo false
  fi
}


# Pushes the current branch with the given name to origin
function push {
  local current_branch_name=$(get_current_branch_name)
  if [ "$(has_tracking_branch "$current_branch_name")" == true ]; then
    if [ "$(needs_push "$current_branch_name")" == true ]; then
      run_command "git push"
    fi
  else
    run_command "git push -u origin $current_branch_name"
  fi
}


# Pushes tags to the remote
function push_tags {
  run_command "git push --tags"
}


# Pushes tags to the remote
function undo_steps_for_push {
  echo "pop_to_next_checkout"
}


function pop_to_next_checkout {
  while [ "$(has_lines "$undo_steps_file")" = true ]; do
    if [[ "$(peek_line "$undo_steps_file")" =~ ^checkout ]]; then
      prepend_to_file "pop_to_next_checkout" "$undo_steps_file"
      break
    else
      remove_line "$undo_steps_file"
    fi
  done
}
