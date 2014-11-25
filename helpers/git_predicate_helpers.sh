#!/bin/bash

# Helper methods that return true or false


# Returns true if the repository has a branch with the given name
function has_branch {
  local branch_name=$1
  if [ "$(git branch | tr -d '* ' | grep -c "^$branch_name\$")" = 0 ]; then
    echo false
  else
    echo true
  fi
}


# Returns true if there are conflicts
function has_conflicts {
  if [ "$(git status | grep -c 'Unmerged paths')" == 0 ]; then
    echo false
  else
    echo true
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


# Determines whether the given branch has a remote tracking branch.
function has_tracking_branch {
  local branch_name=$1
  if [ "$(git branch -r | tr -d ' ' | grep -c "^origin\/$branch_name\$")" == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Determines whether the given branch is ahead of main
function is_ahead_of_main {
  local branch_name=$1
  if [ "$(git log --oneline "$main_branch_name..$branch_name" | wc -l | tr -d ' ')" == 0 ]; then
    echo false
  else
    echo true
  fi
}


# Returns true if the current branch is a feature branch
function is_feature_branch {
  local branch_name=$1
  if [ "$branch_name" == "$main_branch_name" -o "$(echo "$non_feature_branch_names" | tr ',' '\n' | grep -c "$branch_name")" == 1 ]; then
    echo false
  else
    echo true
  fi
}


# Returns whether the current branch has local updates
# that haven't been pushed to the remote yet.
# Assumes the current branch has a tracking branch
function needs_pushing {
  if [ "$(git status | grep -c "Your branch is ahead of")" != 0 ]; then
    echo true
  else
    echo false
  fi
}


# Determines whether the current branch has a rebase in progress
function rebase_in_progress {
  if [ "$(git status | grep -c "rebase in progress")" == 1 ]; then
    echo true
  else
    echo false
  fi
}
