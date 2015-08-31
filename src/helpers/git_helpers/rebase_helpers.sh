#!/usr/bin/env bash


# Abort a rebase
function abort_rebase {
  run_command "git rebase --abort"
}


# Continues rebase if one is in progress
function continue_rebase {
  if [ "$(rebase_in_progress)" == true ]; then
    run_command "git rebase --continue"
  fi
}


# Rebases the given branch into the current branch
function rebase {
  local branch_name=$1
  run_command "git rebase $branch_name"
}


# Determines whether the current branch has a rebase in progress
function rebase_in_progress {
  if [ "$(git status | grep -c "rebase in progress")" == 1 ]; then
    echo true
  else
    echo false
  fi
}
