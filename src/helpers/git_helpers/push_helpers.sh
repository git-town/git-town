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
      run_git_command "git push"
    fi
  else
    run_git_command "git push -u origin $current_branch_name"
  fi
}


# Pushes tags to the remote
function push_tags {
  run_git_command "git push --tags"
}


function undo_steps_for_push {
  echo "skip_current_branch_steps $UNDO_STEPS_FILE"
}
