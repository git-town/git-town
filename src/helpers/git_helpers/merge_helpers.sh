#!/usr/bin/env bash


# Abort a merge
function abort_merge {
  run_command "git merge --abort"
}


# Continues merge if one is in progress
function continue_merge {
  if [ "$(has_open_changes)" == true ]; then
    run_command "git commit --no-edit"
  fi
}


# Merges the given branch into the current branch
function merge {
  local branch_name=$1
  run_command "git merge --no-edit $branch_name"
}


# Squash merges the given branch into the current branch
function squash_merge {
  local branch_name=$1
  local commit_message=$2
  run_command "git merge --squash $branch_name"
  if [ "$commit_message" == "" ]; then
    run_command "git commit -a"
  else
    run_command "git commit -a -m '$commit_message'"
  fi
  if [ $? != 0 ]; then error_empty_commit; fi
}


function undo_steps_for_merge {
  echo "reset_to_sha $(current_sha) hard"
}
