#!/usr/bin/env bash


# Exit if the current branch does not have shippable changes
function ensure_has_shippable_changes {
  local current_branch_name=$(get_current_branch_name)
  if [ "$(has_shippable_changes "$current_branch_name")" == false ]; then
    undo_command

    echo_error_header
    echo_error "The branch '$current_branch_name' has no shippable changes."
    exit_with_error
  fi
}


# Determines whether the given branch has shippable changes
function has_shippable_changes {
  local branch_name=$1
  if [ "$(git diff --quiet "$main_branch_name..$branch_name" ; echo $?)" == 0 ]; then
    echo false
  else
    echo true
  fi
}
