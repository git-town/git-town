#!/usr/bin/env bash


# Commits all open changes into the current branch
function commit_open_changes {
  run_git_command "git add -A"
  run_git_command "git commit -m 'WIP on $(get_current_branch_name)'"
}


# Discard open changes
function discard_open_changes {
  run_git_command 'git reset --hard'
}


# Exists if there are uncommitted changes
function ensure_no_open_changes {
  if [ "$(has_open_changes)" == true ]; then
    echo_error_header
    echo_error "$*"
    exit_with_error newline
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
# and optionally suppress output
function restore_open_changes {
  local silent=$1
  run_git_command "git stash pop" "$silent"
}


# Stashes uncommitted changes
# and optionally suppress output
function stash_open_changes {
  local silent=$1
  run_git_command "git stash -u" "$silent"
}


function undo_steps_for_commit_open_changes {
  local branch=$(get_current_branch_name)
  local sha=$(sha_of_branch "$branch")
  echo "reset_to_sha $sha"
  if [ "$(has_tracking_branch "$branch")" = true ]; then
    echo "push_branch $branch force"
  fi
}


function undo_steps_for_restore_open_changes {
  echo "stash_open_changes"
}


function undo_steps_for_stash_open_changes {
  echo "restore_open_changes"
}
