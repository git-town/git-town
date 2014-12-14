#!/bin/bash


# Commits all open changes into the current branch
function commit_open_changes {
  if [ "$initial_open_changes" = true ]; then
    run_command "git add -A"
    run_command "git commit -m 'WIP on $(get_current_branch_name)'"
  fi
}


# Discard open changes
function discard_open_changes {
  run_command 'git reset --hard'
}


# Exists if there are uncommitted changes
function ensure_no_open_changes {
  if [ "$(has_open_changes)" == true ]; then
    error_has_open_changes

    echo_error_header
    echo_error "$*"
    exit_with_error
  fi
}


# Determines whether there are open changes in Git.
function has_open_changes {
  if [ "$(git status --porcelain | wc -l | tr -d ' ')" == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Unstashes changes
function restore_open_changes {
  run_command "git stash pop"
}


# Stashes uncommitted changes
function stash_open_changes {
  run_command "git stash -u"
}


function undo_steps_for_commit_open_changes {
  local branch=$(get_current_branch_name)
  local sha=$(sha_of_branch "$branch")
  echo "reset_to_sha $sha"
  if [ "$(has_tracking_branch "$branch")" = true ]; then
    echo "push_branch $branch force"
  fi
}


function undo_steps_for_stash_open_changes {
  echo "restore_open_changes"
}
